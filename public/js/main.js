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

function submitAccessForm(form) {
    var mail = form.email.value;
    var username = form.username.value;
    var password = form.password.value;
    var token = form.csrf_token.value;

    var xhr = new XMLHttpRequest();
    var data = 'username=' + username;
    data += '&password=' + password;

    if (form.submit.value === 'Login') {
        xhr.open('POST', '/login', true);
        xhr.onload = function() {
            response = JSON.parse(xhr.responseText);
            
            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                showAlert('You have been successfully logged in!', 'success');
                window.setTimeout(function() {
                    window.location.href = window.location.href;
                }, 3000);
            }
        }
    } else {
        data += '&email=' + mail;
        xhr.open('POST', '/signup', true);
        xhr.onload = function() {
            response = JSON.parse(xhr.responseText);
              
            if (response.error) {
                showAlert(response.message, 'error');
            } else {
                showAlert('You have been successfully signed up!', 'success');
                window.setTimeout(function() {
                    window.location.href = window.location.href;
                }, 3000);
            }
        }
    }
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.setRequestHeader('X-CSRF-Token', token);
    xhr.send(data);
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