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

		Bookmark.prototype.validate = function () {
			var errors = [];

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
      this.token = this.form.csrf_token.value;
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

    App.prototype.addBookmark = function () {
			var bookmark,
			errors,
			data = "",
			tags = [],
			url = "",
			bookmarkUUID,
			i;

			app.getFormValues();

		  bookmark = new Bookmark(null, app.title, app.url, null, null);

			errors = bookmark.validate()
				
      if (errors.length > 0) {
          app.showAlert(errors.join(' '), 'error');
          return false;
      }

      data += '&url=' + bookmark.url;

  		url = '/bookmark/new';

      app.AJAXRequest('POST', 
											url, 
											data, 
											function(response) {
        								if (response.error) {
            							app.showAlert(response.message, 'error');
        								} else {
            							app.showAlert('Bookmark updated successfully.', 'success');
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
												app.renderTags();
                    }
                },
                document.getElementById('csrf_token').value
            );
        }
    }

		App.prototype.listByTag = function () {
			console.log(this);
		}

		/* ==== EVENTS ==== */

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
        // ['load-more-button', loadMore, 'onclick'],
        ['access-form', submitAccessForm, 'onsubmit'],
        ['bookmark-add', app.addBookmark, 'onsubmit'],
        ['no-account', accessFormChangeMode, 'onclick'],
        ['search-form', app.searchBookmarks, 'onsubmit'],
        ['url', app.toggleBookmarkForm, 'onclick'],
			],
			i;

			for (i = 0; i < events.length; i++) {
				setEvent(events[i][0], events[i][1], events[i][2]);
			}

      setKlassEvent('clickable', app.listByTag);
      setKlassEvent('bookmark-delete', app.deleteBookmark);
      setKlassEvent('bookmark-edit', app.editBookmark);
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
