/*global document:false,window:false,XMLHttpRequest:false */
(function() {
    "use strict";

		var Bookmarks = [],
		Tags,
		App;

		/* ==== TAG ==== */

		var Tag = function() {
		};

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
            ulNode,
            tagList,
            newTags,
            i,
            tag,
            tagCount;

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

		/* ==== BOOKMARK ==== */

		var Bookmark = function (uuid, title, url, tags, date) {
      // bookmarkHtml += '<span class="bookmark-tag">' + this.tags[i].trim().toLowerCase() + '</span>';
			this.uuid = uuid;
      this.title = title;
      this.url = url;
      this.tags = tags;
      this.date = date;

			if (this.date === null) {
				this.date = 'Just now';
			}
		};

		var b = new Bookmark("Duude");

		Bookmark.prototype.validate = function () {
			var errors = [];

      if (this.title.length < 1) {
          errors.push('Title cannot be blank.');
      }

      if (this.url.length < 5 || !(this.url.indexOf('http://') !== -1 || this.url.indexOf('https://') !== -1)) {
          errors.push('Invalid url.');
      }

			return errors;
		};


    Bookmark.prototype.render = function() {
				var i = 0,
				bookmarkHtml;

        bookmarkHtml = '<article id="bookmark_' + this.uuid + '">' +
            '<div class="bookmark-actions">' +
            '<a href="#" class="bookmark-edit"><span class="ion-levels"></span></a>' +
            '<a href="#" class="bookmark-delete"><span class="ion-trash-b"></span></a>' +
            '</div>' +
            '<h3><a href="' + this.url + '" target="_blank">' + this.title + '</a></h3>' +
            '<div class="bookmark-url"><span class="ion-link bookmark-icon"></span> ' + this.url + '</div> ' +
            '<div class="bookmark-date"><span class="ion-clock bookmark-icon"></span> ' + this.date + '</div>' +
            '<div class="bookmark-tags"><span class="ion-ios7-pricetag bookmark-icon"></span>';

        if (tags.length === 0) {
            bookmarkHtml += '<span class="bookmark-tag">No tags</span>';
        } else {
            for (i = 0; i < this.tags.length; ++i) {
                bookmarkHtml += '<span class="bookmark-tag">' + this.tags[i] + '</span>';
            }
        }

        bookmarkHtml += '</div></article>';

        return bookmarkHtml;
    }


    function editBookmarkResponse(response, oldTags, tags, data, bookmarkId, date, form) {
        var oldTagsArray,
            tagsArray,
            tagsToAdd,
            tagsToDelete,
            currBk;

        if (response.error) {
            app.showAlert(response.message, 'error');
        } else {
            app.showAlert('Bookmark updated successfully.', 'success');
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

            app.closeForm();

            var viewportOffset = currBk.getBoundingClientRect();

            window.scrollTo(0, viewportOffset.top);

            resetEvents();
        }
    }


    function updateBookmarks(response, list, tag) {
        var data,
            i = 0;

        if (response.error) {
            app.showAlert(response.message, 'error');
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
                    if (null !== document.getElementById('load-more')) {
                        document.getElementById('load-more').style.display = 'none';
                    }
                }
            } else {
                app.showAlert('There are no bookmarks for tag "' + tag + '"', 'info')
            }
            resetEvents();
        }
    }

    function getBookmarksForTag() {
        var form = document.getElementById('bookmark-add'),
            token = form.csrf_token.value,
            list = document.getElementById('list-bookmarks'),
            i = 0,
            tag = this.innerHTML.split(" ")[0];

        app.AJAXRequest('GET', '/tag/' + tag + '/0', '', function(response) {
            updateBookmarks(response, list, tag);
        }, token);
    }

    function searchBookmarksResponse(response, list, query) {
        var i = 0,
            data;

        if (response.error) {
            app.showAlert(response.message, 'error');
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
                    if (null !== document.getElementById('load-more')) {
                        document.getElementById('load-more').style.display = 'none';
                    }
                }
            } else {
                app.showAlert('There are no bookmarks for "' + query + '"', 'info')
            }
        }
    }

		/* ==== APP ==== */

		var heightCallback = function() {
        var docHeight = document.body.scrollHeight;

        if (null !== document.getElementById('left-side')) {
            document.getElementById('left-side').style.height = docHeight + 'px';
        }

        if (null !== document.getElementById('left-side')) {
            document.getElementById('left-side').style.minHeight = docHeight + 'px';
        }
    };

    window.onload = heightCallback;
    window.onresize = heightCallback;

		var App = function(payload) {
			var i,
			bookmark,
			currentSearch,
			currentTags,
			currentPage;

			this.editMode = false;

			this.getFormValues();

			for (i = 0; i < payload.length; i++) {
				bookmark = new Bookmark(payload[i].uuid, payload[i].title, payload[i].url, payload[i].tags, payload[i].date);
				Bookmarks.push(bookmark);
			}
		};

		App.prototype.renderBookmarks = function () {
			var i,
			bookmark;

			document.getElementById('list-bookmarks').innerHTML = "";

			if (Bookmarks.length < 1) {
				document.getElementById('list-bookmarks').innerHTML = '<article class="empty"><h3><span class="ion-ios7-glasses-outline"></span></h3><p>There aren\'t any bookmarks yet.</p></article>';
				return;
			}

			for (i = 0; i < Bookmarks.length; i++) {
			  document.getElementById('list-bookmarks').innerHTML += Bookmarks[i].render();
			}
		};

		App.prototype.getFormValues = function () {
      this.form = document.getElementById('bookmark-add'),
      this.title = this.form.title.value,
      this.url = this.form.url.value,
      this.tags = this.form.tags.value,
      this.token = this.form.csrf_token.value;
		};

    App.prototype.toggleBookmarkForm = function () {
			app.form.title.parentNode.classList.toggle('hidden');
			app.form.tags.parentNode.classList.toggle('hidden');
    };

    App.prototype.editBookmark = function () {
        var uuidNode = this.parentNode.parentNode,
						i = 0;

        app.toggleBookmarkForm();

				var bookmark;

				var bookmarkUUID = uuidNode.id.substring(uuidNode.id.indexOf('_') + 1);

				for (i = 0; i < Bookmarks.length; i++) {
					if ( Bookmarks[i].uuid === bookmarkUUID) {
						bookmark = Bookmarks[i];
					} 
				}

        app.form.submit.value = 'Edit bookmark';
        app.form.tags.value = bookmark.tags;
        app.form.bookmark_id.value = bookmark.uuid
        app.form.old_tags.value = bookmark.tags;
        app.form.title.value = bookmark.title;
        app.form.url.value = bookmark.url;
        app.form.bookmark_id.value = bookmark.uuid;
        app.form.bookmark_date.value = bookmark.date;

        document.getElementById('toggle_edit_form').className = 'button-action';

				app.editMode = true;

        window.scrollTo(0, 0);
    }

		App.prototype.closeForm = function () {
      this.form.submit.value = 'Add bookmark';
      this.form.tags.value = '';
      this.form.title.value = '';
      this.form.url.value = '';

			app.toggleBookmarkForm();
		};

    App.prototype.AJAXRequest = function (method, url, data, callback, token) {
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

    App.prototype.refresh = function () {
        window.setTimeout(function() {
            window.location.href = window.location.href;
        }, 3000);
    }

    App.prototype.escapeHTMLEntities = function (str) {
        return str.replace(/[&<>]/g, function(entity) {
            return {
                '&': '&amp;',
                '<': '&lt;',
                '>': '&gt;'
            } || entity;
        });
    }

    App.prototype.showAlert = function (msg, htmlClass) {
        var alert = document.getElementById('alert');
        alert.className = htmlClass;
        alert.innerHTML = msg;
        alert.style.display = 'block';

        window.setTimeout(function() {
            alert.style.display = 'none';
        }, 2000);
    }

    App.prototype.searchBookmarks = function () {
        var form = document.getElementById('bookmark-add'),
            token = form.csrf_token.value,
            list = document.getElementById('list-bookmarks'),
            query = document.getElementById('search-form').search_query.value;

        app.AJAXRequest( 'POST', '/search/0', 'query=' + query, function(response) { searchBookmarksResponse(response, list, query); }, token );

        return false;
    }

		App.prototype.stringSplitTrim = function (string) {
			var i,
			temp,
			result = [];

			temp = string.replace(/(\,+)/g, ",").replace(/$\,|^\,/g, "").trim().split(",");

			for (i = 0; i < temp.length; i++) {
				result.push(temp[i].trim());
			}

			return result;
		}

    App.prototype.addBookmark = function () {
			var bookmark,
			errors,
			data = "",
			tags = [],
			url = "",
			bookmarkUUID,
			i;

			app.getFormValues();

			tags = app.stringSplitTrim(app.tags);

			if (app.editMode === true) {
				bookmarkUUID = app.form.bookmark_id.value;

				for (i = 0; i < Bookmarks.length; i++) {
					if ( Bookmarks[i].uuid === bookmarkUUID) {
						bookmark = Bookmarks[i];
					} 
				}

			  bookmark.title = app.title;
				bookmark.url = app.url;
				bookmark.tags = tags;
				bookmark.date = 'Just now';
			} else {
			  bookmark = new Bookmark(null, app.title, app.url, tags, null);
			}

			errors = bookmark.validate()
				
      if (errors.length > 0) {
          app.showAlert(errors.join(' '), 'error');
          return false;
      }

      data += 'title=' + bookmark.title;
      data += '&url=' + bookmark.url;
      data += '&tags=' + bookmark.tags;

			if (app.editMode === false) {
				url = '/bookmark/new';
			} else {
				url = '/bookmark/update/' + bookmark.uuid;
			}

      app.AJAXRequest('POST', 
											url, 
											data, 
											function(response) {
        								if (response.error) {
            							app.showAlert(response.message, 'error');
        								} else {
													if (app.editMode === false) {
													  bookmark.uuid = response.message;
													  Bookmarks.unshift(bookmark);
            							  app.showAlert('Bookmark added successfully.', 'success');
														app.editMode = false;
													} else {
            							  app.showAlert('Bookmark updated successfully.', 'success');
													}

													app.closeForm();
													app.renderBookmarks();
													resetEvents();
        								}
			 								}, 
											app.token);

     return false;
    }

		var app = new App(payload);

		app.getFormValues();
		app.renderBookmarks();

    App.prototype.deleteBookmark = function () {
        var elem = this.parentNode.parentNode,
            id = elem.id.split("_")[1];

        if (confirm("Are you sure you want to delete that?")) {
            app.AJAXRequest(
                'DELETE',
                '/bookmark/delete/' + id,
                '',
                function(response) {
                    if (response.error) {
                        app.showAlert(response.message, 'error');
                    } else {
                        app.showAlert('Bookmark deleted successfully.', 'success');
                        elem.style.display = 'none';
                        updateTags(getTagsFromBookmark(elem), true);
                    }
                },
                document.getElementById('csrf_token').value
            );
        }
    }

    function setEvent(element, func, prop) {
        if (null !== document.getElementById(element)) {
            document.getElementById(element)[prop] = func;
        }
    }

    function setKlassEvent(klass, func) {
        var nodes = document.getElementsByClassName(klass),
				i;

        if (typeof(nodes) !== 'undefined' && nodes.length > 0) {
            for (i = 0; i < nodes.length; ++i) {
                nodes[i].onclick = func;
            }
        }
    }

    function resetEvents() {
			var events = [
        // ['browseALL', browseAll, 'onclick'],
        ['no-account', accessFormChangeMode, 'onclick'],
        ['url', app.toggleBookmarkForm, 'onclick'],
        ['toggle_edit_form', app.closeForm, 'onclick'],
        ['access-form', submitAccessForm, 'onsubmit'],
        ['bookmark-add', app.addBookmark, 'onsubmit'],
        ['search-form', app.searchBookmarks, 'onsubmit'],
        // ['load-more-button', loadMore, 'onclick'],
			],
			i;

			for (i = 0; i < events.length; i++) {
				setEvent(events[i][0], events[i][1], events[i][2]);
			}

      setKlassEvent('bookmark-edit', app.editBookmark);
      setKlassEvent('bookmark-delete', app.deleteBookmark);
      setKlassEvent('clickable', getBookmarksForTag);
    }

    window.addEventListener('load', resetEvents(), false);

		/* ==== LOGIN ==== */
		
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

    function submitAccessForm() {
        var form = document.getElementById('access-form'),
            mail = form.email.value,
            username = form.username.value,
            password = form.password.value,
            token = form.csrf_token.value,
            data = 'username=' + username;

        data += '&password=' + password;

        if (form.submit.value !== 'Login') {
            data += '&email=' + mail;
        }

        app.AJAXRequest(
            'POST',
            form.submit.value === 'Login' ? '/login' : '/signup',
            data,
            form.submit.value === 'Login' ? function(response) {
                if (response.error) {
                    app.showAlert(response.message, 'error');
                } else {
                    app.showAlert('You have been successfully logged in!', 'success');
                    app.refresh();
                }
            } : function(response) {
                if (response.error) {
                    app.showAlert(response.message, 'error');
                } else {
                    app.showAlert('You have been successfully signed up!', 'success');
                    app.refresh();
                }
            },
            token
        );

        return false;
    }
}());
