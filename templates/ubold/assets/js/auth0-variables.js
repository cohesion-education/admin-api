var AUTH0_CLIENT_ID='DBfgngEpPVlRawcfFWme3gxJ6WNwBEl6';
var AUTH0_DOMAIN='cohesioned.auth0.com';
var AUTH0_CALLBACK_URL=window.location.href + 'callback';
var RETURN_TO_URL='http://localhost:3000/logout'
var AUTH0_LOGOUT_URL='https://'+AUTH0_DOMAIN+'/v2/logout?returnTo='+encodeURI(RETURN_TO_URL)+'&client_id='+AUTH0_CLIENT_ID;
