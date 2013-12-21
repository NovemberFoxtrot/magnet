/*global document:false,window:false,XMLHttpRequest:false */
(function() {
    "use strict";

		var Bookmarks = [],
		Tags,
		App;

		/* ==== TAG ==== */

		var Tag = function(title, count) {
			this.title = title;
			this.count = count;
		};

		Tag.prototype.render = function () {
			return '<li class="clickable">' + this.title + '<span class="tag-count">(' + this.count + ')</span></li>';
		}

    function recountTags() {
			var i,
			j,
			stats = {};

			for (i = 0; i < Bookmarks.length; i++) {
				for (j = 0; j < Bookmarks[i].tags.length; j++ ) {
					if (typeof(stats[Bookmarks[i].tags[j]]) === "undefined") {
						stats[Bookmarks[i].tags[j]] = 1;
					} else {
						stats[Bookmarks[i].tags[j]] += 1;
					}
				}
			}

			return stats;
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
		  // document.getElementById('toggle_edit_form').className.remove('button-action');
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

        app.AJAXRequest('POST', '/search/0', 'query=' + query, function(response) { searchBookmarksResponse(response, list, query); }, token);

        return false;
    }

		App.prototype.stringSplitTrim = function (string) {
			var i,
			temp,
			result = [];

			// replace multiple ,s with a singe , // remove , from beginning and end 
			temp = string.replace(/(\,+)/g, ",").replace(/$\,|^\,/g, "").trim().split(",");

			for (i = 0; i < temp.length; i++) {
				result.push(temp[i].trim().toLowerCase());
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

		App.prototype.renderTags = function () {
			var tagCounts = recountTags(),
			i;

			for (i = 0; i < tagCounts.length; i++) {
				console.log(tagCounts[i]);
			}
		}

		var app = new App(payload);

		app.getFormValues();
		app.renderBookmarks();
		app.renderTags();

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
        ['submit-add-bookmark', app.openForm, 'onclick'],
        ['search-form', app.searchBookmarks, 'onsubmit'],
        // ['load-more-button', loadMore, 'onclick'],
			],
			i;

			for (i = 0; i < events.length; i++) {
				setEvent(events[i][0], events[i][1], events[i][2]);
			}

      setKlassEvent('bookmark-edit', app.editBookmark);
      setKlassEvent('bookmark-delete', app.deleteBookmark);
      // setKlassEvent('clickable', getBookmarksForTag);
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
