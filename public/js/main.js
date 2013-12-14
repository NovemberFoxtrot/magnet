(function () {
"use strict";

window.addEventListener('load', theLockAndLoad(), false);

function lock_and_load(element, func) {
  if (null !== document.getElementById(element)) {
  	document.getElementById(element).onclick = func;
  }
}

function theLockAndLoad() {
	console.log("lock n' load");

 	lock_and_load('browseAll', browseAll);
  lock_and_load('access-form', submitAccessForm);
  lock_and_load('no-account', accessFormChangeMode);
  lock_and_load('bookmark-add', submitNewBookmark);
  lock_and_load('url', toggleBookmarkForm); // (true) true
  lock_and_load('toggle_edit_form', closeEditBookmarkForm); // (this.parentNode.parentNode)

	// need to loop over class?
  lock_and_load('bookmark-edit', openEditBookmarkForm); // this.parentNode.parentNode
  lock_and_load('bookmark-delete', deleteBookmark); // ('{{Id}}', this.parentNode.parentNode)

  lock_and_load('load-more-button', loadMore); // (1)
  lock_and_load('clickable', getBookmarksForTag); // ('{{Name}}')
}

var heightCallback = function() {
    var docHeight = document.body.scrollHeight;
    document.getElementById('left-side').style.height = docHeight + 'px';
    document.getElementById('left-side').style.minHeight = docHeight + 'px';
};

window.onload = heightCallback;
window.onresize = heightCallback;

function accessFormChangeMode() {
    var submit = document.getElementById('submit-button'),
    modeChanger = document.getElementById('no-account'),
    email = document.getElementById('email-field');
    
    if (submit.value === 'Login') {
        email.className = 'form-field';
        submit.value = 'Sign up';
        modeChanger.value = 'I have an account';
    } else {
        email.className = 'form-field hidden';
        submit.value = 'Login';
        modeChanger.value = 'I don\'t have an account';
    }
}

function AJAXRequest(method, url, data, callback, token) {
    var xhr = new XMLHttpRequest(),
		response;

    xhr.open(method, url, true);

    xhr.onload = function() {
        response = JSON.parse(xhr.responseText);
        callback(response);
    };

    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

    if (token !== undefined) {
        xhr.setRequestHeader('X-CSRF-Token', token);
		}

    xhr.send(data);
}

function refresh() {
    window.setTimeout(function() {
        window.location.href = window.location.href;
    }, 3000);
}

function escapeHTMLEntities(str) {
    return str.replace(/[&<>]/g, function(entity) {
        return {
            '&' : '&amp;',
            '<' : '&lt;',
            '>' : '&gt;'
        } || entity;
    });
}

function submitAccessForm() {
    var form = document.getElementById('access-form');

    var mail = form.email.value;
    var username = form.username.value;
    var password = form.password.value;
    var token = form.csrf_token.value;
    var data = 'username=' + username;

    data += '&password=' + password;

    if (form.submit.value !== 'Login') {
        data += '&email=' + mail;
		}

    AJAXRequest(
        'POST',
        form.submit.value === 'Login' ? '/login' : '/signup',
        data,
        form.submit.value === 'Login' ? function(response) {
            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                showAlert('You have been successfully logged in!', 'success');
                refresh();
            }
        } : function(response) {
            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                showAlert('You have been successfully signed up!', 'success');
                refresh();
            }
        },
        token
    );
}

function submitNewBookmarkResponse(response, tags) {
	var empty,
	lb;

            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                showAlert('Bookmark added successfully.', 'success');
                empty = document.getElementsByClassName('empty');
                if (empty.length > 0) {
                    empty[0].style.display = 'none';
                }
 
                lb = document.getElementById('list-bookmarks');
                lb.innerHTML = renderBookmark(response.message, title.value, url.value, tags.value) + lb.innerHTML;
                updateTags(tags.value);
                title.value = '';
                url.value = '';
                tags.value = '';
                toggleBookmarkForm(false);
            }
	
}

function submitNewBookmark(form) {
    var title = form.title,
        url = form.url,
        tags = form.tags,
        token = form.csrf_token.value,
        data = '',
        errorMessages = [];
        
    if (title.value.length < 1) {
        errorMessages.push('Title cannot be blank.');
    }
    
    if (url.value.length < 5 || 
        !(url.value.indexOf('http://') !== -1 || url.value.indexOf('https://') !== -1)) {
        errorMessages.push('Invalid url.');
    }
    
    if (errorMessages.length > 0) {
        showAlert(errorMessages.join(' '), 'error');
        return;
    }

    data += 'title=' + title.value;
    data += '&url=' + url.value;
    data += '&tags=' + tags.value;

    AJAXRequest('POST', '/bookmark/new', data, function(response) { submitNewBookmarkResponse(response, tags); }, token);
}

function showAlert(msg, htmlClass) {
    var alert = document.getElementById('alert');
    alert.className = htmlClass;
    alert.innerHTML = msg;
    alert.style.display = 'block';
    
    window.setTimeout(function() {
        alert.style.display = 'none';
    }, 2000);
}

function renderBookmark(bkId, title, url, tags, date, forceComplete) {
    var editing = true && !forceComplete,
		bookmarkHtml,
		tagArray,
		i;

    if (date === undefined) {
        date = 'Just now';
        editing = false;
    }

    title = escapeHTMLEntities(title);
    url = escapeHTMLEntities(url);
    tags = escapeHTMLEntities(tags);

	bookmarkHtml = ((!editing) ? '<article id="bookmark_' + bkId + '">' : '') + 
        '<div class="bookmark-actions">' +
		'<a href="#" class="bookmark-edit"><span class="ion-levels"></span></a>' +
		'<a href="#" class="bookmark-delete"><span class="ion-trash-b"></span></a>' +
		'</div>' +
		'<h3><a href="'+ url + '" target="_blank">' + title + '</a></h3>' +
		'<div class="bookmark-url"><span class="ion-link bookmark-icon"></span> ' + url + '</div> ' + 
        '<div class="bookmark-date"><span class="ion-clock bookmark-icon"></span> ' + date + '</div>' +
        '<div class="bookmark-tags"><span class="ion-ios7-pricetag bookmark-icon"></span>';
    if (tags.trim() === '') {
        bookmarkHtml += '<span class="bookmark-tag">No tags</span>';
    } else {
        tagArray = tags.split(',');
        for (i in tagArray) {
            bookmarkHtml += '<span class="bookmark-tag">' + tagArray[i].trim().toLowerCase() + '</span>';
        }
    }

    bookmarkHtml += '</div>' + ((!editing) ? '</article>' : '');
    
    return bookmarkHtml;
}

function getTagsFromBookmark(bookmarkElem) {
    var tagElems = bookmarkElem.getElementsByClassName('bookmark-tag'),
    tags = [],
		tag,
		i;

    for (i = 0; i < tagElems.length; i++) {
        tag = tagElems[i].innerHTML;
        if (tag !== 'No tags') tags.push(tag);
    }
    
    return tags.join(', ');
}

function appendTag(ulNode, tag, tagCount) {
    if (tagCount === undefined) tagCount = 1;
    ulNode.innerHTML += '<li class="clickable">' + tag + ' <span class="tag-count">(' + tagCount + ')</span></li>';
}

function updateTags(tags, deleteTags) {
		var tagArray,
		ulNode;

    if (tags.trim() !== '') {
        tagArray = tagsToArray(tags);
        
        if (deleteTags === undefined) deleteTags = false;
        
        ulNode = document.getElementById('tags').
            getElementsByTagName('ul')[0];
        tagList = ulNode.getElementsByTagName('li');
            
        // If there is only one element and it's not clickable it means
        // that this element is the No tags placeholder.
        if (tagList.length === 1 && tagList[0].className !== 'clickable' && !deleteTags) {
            tagList[0].style.display = 'none';

            for (i in tagArray) {
                appendTag(ulNode, tagArray[i]);
            }

        } else {
            newTags = [];
            
            for (i = 0; i < tagList.length; i++) {
                tag = tagList[i].innerHTML.split(' <span class="tag-count">')[0];
                tagCount = Number(tagList[i].innerHTML.split('(')[1].split(')')[0]);

                if (tagArray.indexOf(tag) !== -1) {
                    tagCount += (deleteTags ? -1 : 1);
                    tagArray[tagArray.indexOf(tag)] = null;
                }
                
                if (tagCount > 0) {
                    newTags.push([tag, tagCount]);
                }
            }
            
            ulNode.innerHTML = '';

            for (i in newTags) {
                appendTag(ulNode, newTags[i][0], newTags[i][1]);
            }

            if (!deleteTags) {
                for (i in tagArray) {
                    if (tagArray[i] !== null) {
                        appendTag(ulNode, tagArray[i]);
                    }
                }
            }
        }
    }
}

function tagsToArray(tags) {
    var tagArray = tags.split(',');
    for (var i = 0; i < tagArray.length; i++) {
        tagArray[i] = tagArray[i].trim().toLowerCase();
    }

    return tagArray;
}

function toggleBookmarkForm(open) {
    var form = document.getElementById('bookmark-add'),
        fields = form.getElementsByClassName('form-field'),
        htmlClass = (open) ? '' : ' hidden';
        
    for (var i in fields) {
        if (i > 0) {
            fields[i].className = 'form-field' + htmlClass;
				}
    }

    form.getElementsByClassName('form-buttons').className = 'form-buttons' + htmlClass;
}

function deleteBookmark(id, elem) {
    if (confirm("Are you sure you want to delete that?")) {
        AJAXRequest(
            'DELETE',
            '/bookmark/delete/' + id,
            '',
            function(response) {
                if (response.error) {
                    showAlert(response.message, 'error');
                } else {
                    showAlert('Bookmark deleted successfully.', 'success');
                    elem.style.display = 'none';
                    updateTags(getTagsFromBookmark(elem), true);
                }
            },
            document.getElementById('csrf_token').value
        );
    }
}

function editBookmarkResponse(response, oldTags, tags, data) {
            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                showAlert('Bookmark updated successfully.', 'success');
                // Update tags
                if (tags.value !== '') {
                    if (oldTags.value === '') {
                        // Add the new tags
                        updateTags(tags.value);
                    } else {
                        // Remove those which where removed and 
                        // add the new ones
                        oldTagsArray = tagsToArray(oldTags.value);
                        tagsArray = tagsToArray(tags.value);
                        tagsToAdd = [];
                        tagsToDelete = [];

                        for (var i = 0; i < tagsArray.length; i++) {
                            if (oldTagsArray.indexOf(tagsArray[i]) === -1) {
                                tagsToAdd.push(tagsArray[i]);
                            } else {
                                oldTagsArray[oldTagsArray.indexOf(tagsArray[i])] = undefined;
                            }
                        }

                        for (i = 0; i < oldTagsArray.length; i++) {
                            if (oldTagsArray[i] !== undefined) {
                                tagsToDelete.push(oldTagsArray[i]);
                            }
                        }

                        updateTags(tagsToAdd.join(', '));
                        updateTags(tagsToDelete.join(', '), true);
                    }
                } else if (oldTags !== '') {
                    // Remove the tags
                    updateTags(oldTags.value, true);
                }
                currBk = document.getElementById('bookmark_' + bookmarkId.value);
                currBk.innerHTML = renderBookmark(bookmarkId.value, title.value, url.value, tags.value, date.value);
                closeEditBookmarkForm(form);
                var viewportOffset = currBk.getBoundingClientRect();
                window.scrollTo(0, viewportOffset.top);
            }
}

function editBookmark(form) {
    var title = form.title,
        url = form.url,
        tags = form.tags,
        token = form.csrf_token.value,
        bookmarkId = form.bookmark_id,
        date = form.bookmark_date,
        oldTags = form.old_tags,
        data = '',
        errorMessages = [];
        
    if (title.value.length < 1) {
        errorMessages.push('Title cannot be blank.');
    }
    
    if (url.value.length < 5 || !(url.value.indexOf('http://') !== -1 || url.value.indexOf('https://') !== -1)) {
        errorMessages.push('Invalid url.');
    }
    
    if (errorMessages.length > 0) {
        showAlert(errorMessages.join(' '), 'error');
        return;
    }

    data += 'title=' + title.value;
    data += '&url=' + url.value;
    data += '&tags=' + tags.value;

    AJAXRequest(
        'POST',
        '/bookmark/update/' + bookmarkId.value,
        data,
        function(response) {
					editBookmarkResponse(response, oldTags, tags, data);
        },
        token
    );
}

function openEditBookmarkForm(bookmark) {
    var form = document.getElementById('bookmark-add'),
		dateElemContent;

    toggleBookmarkForm(true);

    form.onsubmit = function() {
        editBookmark(form);
        return false;
    }

    form.submit.value = 'Edit bookmark';

    form.tags.value = getTagsFromBookmark(bookmark);

    form.bookmark_id.value = bookmark.id.substring(bookmark.id.indexOf('_') + 1);

    form.old_tags.value = form.tags.value;

    form.title.value = bookmark.getElementsByTagName('h3')[0].getElementsByTagName('a')[0].innerHTML;

    form.url.value = bookmark.getElementsByClassName('bookmark-url')[0].innerHTML.split(' ')[3];

    dateElemContent = bookmark.getElementsByClassName('bookmark-date')[0].innerHTML;

    form.bookmark_date.value = dateElemContent.substring(dateElemContent.lastIndexOf('>') + 2);

    document.getElementById('toggle_edit_form').className = 'button-action';

    window.scrollTo(0, 0);
}

function closeEditBookmarkForm(form) {
    form.onsubmit = function() {
        submitNewBookmark(form);
        return false;
    }

    form.submit.value = 'Add bookmark';
    form.tags.value = '';
    form.title.value = '';
    form.url.value = '';

    document.getElementById('toggle_edit_form').className = 'button-action hidden';
}

function updateBookmarks(response, list, tag) {
			var data,
			i = 0;

				if (response.error) {
								showAlert(response.message, 'error');
				} else {
								data = response.data;
								list.className = 'browsing_tag_' + tag;
								if (data.length > 0) {
												list.innerHTML = '';
												for (i = 0; i < data.length; i++) {
																list.innerHTML += renderBookmark(data[i].id,
																								data[i].Title,
																								data[i].Url,
																								data[i].Tags.join(', '),
																								data[i].Date,
																								true);
												}

												document.getElementById('back-index').className = '';

												if (data.length == 50) {
																document.getElementById('load-more').onclick = function() {
																				loadMore(1);
																				return false;
																};
												} else {
																document.getElementById('load-more').style.display = 'none';
												}
								} else {
												showAlert('There are no bookmarks for tag "' + tag + '"', 'info')
								}
				}
}

function getBookmarksForTag(tag) {
    var form = document.getElementById('bookmark-add'),
        token = form.csrf_token.value,
        list = document.getElementById('list-bookmarks'),
        i = 0;

    AJAXRequest(
        'GET',
        '/tag/' + tag + '/0',
        '',
				function(response) {
				  updateBookmarks(response, list, tag)
				},
        token
    );
}

function searchBookmarksResponse(response) {
        var i = 0;

				if (response.error) {
								showAlert(response.message, 'error');
				} else {
								data = response.data;
								list.className = 'searching_' + btoa(query);
								if (data.length > 0) {
												list.innerHTML = '';
												for (i = 0; i < data.length; i++) {
																list.innerHTML += renderBookmark(data[i].id,
																								data[i].Title,
																								data[i].Url,
																								data[i].Tags.join(', '),
																								data[i].Date,
																								true);
												}
                    
												document.getElementById('back-index').className = '';

												if (data.length == 50) {
																document.getElementById('load-more').onclick = function() {
																				loadMore(1);
																				return false;
																};
												} else {
																document.getElementById('load-more').style.display = 'none';
												}
								} else {
												showAlert('There are no bookmarks for tag "' + tag + '"', 'info')
								}
				}
}

function searchBookmarks(query) {
    var form = document.getElementById('bookmark-add'),
        token = form.csrf_token.value,
        list = document.getElementById('list-bookmarks');

    AJAXRequest(
        'POST',
        '/search/0',
        'query=' + query,
        function(response) {
					searchBookmarksResponse(response, list);
        },
        token
    );
}

function browseAllResponse(response, list) {
	var i = 0;

            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                data = response.data;
                if (list.className.indexOf('searching_') !== -1) {
                    document.getElementById('search_query').value = '';
                }
                list.className = '';
                if (data.length > 0) {
                    list.innerHTML = '';
                    for (i = 0; i < data.length; i++) {
                        list.innerHTML += renderBookmark(data[i].id,
                                                        data[i].Title,
                                                        data[i].Url,
                                                        data[i].Tags.join(', '),
                                                        data[i].Date,
                                                        true);
                    }
                    
                    document.getElementById('back-index').className = 'hidden';
                    
                    if (data.length == 50) {
                        document.getElementById('load-more').onclick = function() {
                            loadMore(1);
                            return false;
                        };
                    } else {
                        document.getElementById('load-more').style.display = 'none';
                    }
                } else {
                    showAlert('There are no bookmarks to display.', 'info')
                }
            }
}

function browseAll() {
    var form = document.getElementById('bookmark-add'),
        token = form.csrf_token.value,
        list = document.getElementById('list-bookmarks');

    AJAXRequest('GET', '/bookmarks/0', '', function(response) {browseAllResponse(response, list);}, token);
}

function loadMoreResponse(response, list) {
	var i = 0;
            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                data = response.data;
                console.log(data);
                if (data.length > 0) {
                    for (i = 0; i < data.length; i++) {
                        list.innerHTML += renderBookmark(data[i].id,
                                                        data[i].Title,
                                                        data[i].Url,
                                                        data[i].Tags.join(', '),
                                                        data[i].Date,
                                                        true);
                    }
                    
                    if (data.length == 50) {
                        document.getElementById('load-more').onclick = function() {
                            loadMore(page + 1);
                            return false;
                        };
                    } else {
                        document.getElementById('load-more').style.display = 'none';
                    }
                } else {
                    showAlert('There are no bookmarks to display.', 'info')
                }
            }
	
}

function loadMore(page) {
    var form = document.getElementById('bookmark-add'),
        token = form.csrf_token.value,
        list = document.getElementById('list-bookmarks'),
        method,
        queryData,
        requestUrl,
        i = 0;
        
    if (list.className.indexOf('browsing_tag_') !== -1) {
        method = 'GET';
        requestUrl = '/tag/' + list.className.substring(list.className.indexOf('tag_') + 4) + '/' + page;
        queryData = '';
    } else if (list.className.indexOf('searching_') !== -1) {
        method = 'POST';
        requestUrl = '/search/' + page;
        queryData = 'query=' + atob(list.className.substring(list.className.indexOf('_') + 1));
    } else {
        method = 'GET';
        requestUrl = '/bookmarks/' + page;
        queryData = '';
    }
        
    AJAXRequest(method, requestUrl, queryData, function(response) {loadMoreResponse(response, list); }, token);
}
}());
