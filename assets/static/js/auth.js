function attachSignin(element) {
  auth2.attachClickHandler(element, {}, function(googleUser) {
    var form = document.getElementById('login')
    var email = document.getElementById('email');
    var token = document.getElementById('token');

    token.value = googleUser.getAuthResponse().id_token;
    email.value = googleUser.getBasicProfile().getEmail();

    error.textContent = '';
    form.submit();
  }, function(error) {
    console.error(error);
  });
}

var startApp = function() {
  gapi.load('auth2', function(){
    // Retrieve the singleton for the GoogleAuth library and set up the client.
    auth2 = gapi.auth2.init({
      client_id: '420546501001-v35bges8923km4s9r9p3tet8m42ibj5m.apps.googleusercontent.com',
      cookiepolicy: 'single_host_origin',
      // Request scopes in addition to 'profile' and 'email'
      //scope: 'additional_scope'
    });
    attachSignin(document.getElementById('google-signin'));
  });
};
