function accessFormChangeMode() {
    var submit = document.getElementById('submit-button');
    var modeChanger = document.getElementById('no-account');
    var email = document.getElementById('email-field');
    
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
    var xhr = new XMLHttpRequest();
    xhr.open(method, url, true);
    xhr.onload = function() {
        response = JSON.parse(xhr.responseText);
        callback(response);
    };
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    if (token !== undefined)
        xhr.setRequestHeader('X-CSRF-Token', token);
    xhr.send(data);
}

function refresh() {
    window.setTimeout(function() {
        window.location.href = window.location.href;
    }, 3000);
}

function submitAccessForm(form) {
    var mail = form.email.value;
    var username = form.username.value;
    var password = form.password.value;
    var token = form.csrf_token.value;
    var data = 'username=' + username;
    data += '&password=' + password;
    if (form.submit.value !== 'Login')
        data += '&email=' + mail;

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
    // TODO strip HTML tags
    data += '&tags=' + tags.value;

    AJAXRequest(
        'POST',
        '/bookmark/new',
        data,
        function(response) {
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
        },
        token
    );
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

function renderBookmark(bkId, title, url, tags, date) {
    if (date === undefined) {
        date = 'Just now';
    }
    // TODO delete and edit buttons
	bookmarkHtml = '<article>' + 
        '<div class="bookmark-actions">' +
		'<a href="#" class="bookmark-edit" onclick="editBookmark(\'' + bkId + '\', this.parentNode.parentNode);"><span class="ion-levels"></span></a>' +
		'<a href="#" class="bookmark-delete" onclick="deleteBookmark(\'' + bkId + '\', this.parentNode.parentNode);"><span class="ion-trash-b"></span></a>' +
		'</div>' +
		'<h3><a href="'+ url + '">' + title + '</a></h3>' +
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
    bookmarkHtml += '</div></article>';
    
    return bookmarkHtml;
}

function getTagsFromBookmark(bookmarkElem) {
    tagElems = bookmarkElem.getElementsByClassName('bookmark-tag');
    tags = [];
    for (i = 0; i < tagElems.length; i++) {
        tag = tagElems[i].innerHTML;
        if (tag !== 'No tags') tags.push(tag);
    }
    
    return tags.join(', ');
}

function appendTag(ulNode, tag, tagCount) {
    if (tagCount === undefined) tagCount = 1;
    ulNode.innerHTML += '<li class="clickable" onclick="getBookmarksForTag(\'' + 
        tag + '\');">' + tag + ' <span class="tag-count">(' + tagCount + ')</span></li>';
}

function updateTags(tags, deleteTags) {
    if (tags.trim() !== '') {
        tagArray = tags.split(',');
        for (i in tagArray) {
            tagArray[i] = tagArray[i].trim().toLowerCase();
        }
        
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
                tagCount = Number(tagList[i].innerHTML.split('(')[1].
                        split(')')[0]);
                if (tagArray.indexOf(tag) !== -1) {
                    tagCount += (deleteTags ? -1 : 1);
                    tagArray[tagArray.indexOf(tag)] = null;
                }
                
                if (tagCount > 0) {
                    newTags.push([tag, tagCount]);
                    console.log(newTags);
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

function toggleBookmarkForm(open) {
    var form = document.getElementById('bookmark-add'),
        fields = form.getElementsByClassName('form-field'),
        htmlClass = (open) ? '' : ' hidden';
        
    for (var i in fields) {
        if (i > 0)
            fields[i].className = 'form-field' + htmlClass;
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