var AUTH0_CLIENT_ID='DBfgngEpPVlRawcfFWme3gxJ6WNwBEl6';
var AUTH0_DOMAIN='cohesioned.auth0.com';
var base = location.protocol+'//'+location.hostname+(location.port ? ':'+location.port: '');
var AUTH0_CALLBACK_URL=base + '/callback';
var RETURN_TO_URL=base + '/logout'
var AUTH0_LOGOUT_URL='https://'+AUTH0_DOMAIN+'/v2/logout?returnTo='+encodeURI(RETURN_TO_URL)+'&client_id='+AUTH0_CLIENT_ID;
