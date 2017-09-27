function onSignIn(googleUser) {
	var form = document.getElementById('login')
  var user = document.getElementById('user');
  var token = document.getElementById('token');

  token.value = googleUser.getAuthResponse().id_token;
  user.value = googleUser.getBasicProfile().getEmail();
  form.submit();
}


function attachSignin(element) {
	console.log(element.id);
  auth2.attachClickHandler(element, {}, onSignIn, function(error) {
    alert(JSON.stringify(error, undefined, 2));
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
