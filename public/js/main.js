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
        data = '';

    data += 'title=' + title.value;
    data += '&url=' + url.value;
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

function updateTags(tags) {
    if (tags.trim() !== '') {
        tagArray = tags.split(',');
        for (i in tagArray) {
            tagArray[i] = tagArray[i].trim().toLowerCase();
        }
    }

    // Look for the tags in the tag list
    // Update their count
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