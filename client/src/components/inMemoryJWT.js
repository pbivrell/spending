import axios from 'axios';

const inMemoryJWTManager = () => {
    let logoutEventName = 'ra-logout';
    let Token = null;
    let refreshTimeOutId;
    let refreshing = false;
    let isRefreshing = null;
    let tokenResponse = null;

    // This listener allows to disconnect another session of react-admin started in another tab
    window.addEventListener('storage', (event) => {
        if (event.key === logoutEventName) {
            Token = null;
	    console.log("SET TOKEN TO NULL");
        }
    });

    // This countdown feature is used to renew the JWT in a way that is transparent to the user.
    // before it's no longer valid
    const refreshToken = (delay) => {

        refreshTimeOutId = window.setTimeout(
            getRefreshedToken,
            delay * 1000 - 5000
        ); // Validity period of the token in seconds, minus 5 seconds
    };

    const abordRefreshToken = () => {
        if (refreshTimeOutId) {
            window.clearTimeout(refreshTimeOutId);
        }
    };
 
    const waitForTokenRefresh = () => {
        if (!refreshing) {
            return Promise.resolve();
        }
        return isRefreshing.then(() => {
            isRefreshing = null;
            return true;
        });
    }

    // The method makes a call to the refresh-token endpoint
    // If there is a valid cookie, the endpoint will return a fresh jwt.
    const getRefreshedToken = () => {
	if (refreshing)  {
		return;
	}

	refreshing = true;

	isRefreshing = axios.get('http://localhost:8080/api/v1/refresh',   {
		withCredentials: true
	})
	.then(function(response){
		setToken(response.data);
		refreshing = false;
	})
	.catch(function(error) {
		console.log(error);
		refreshing = false;
	})
	return isRefreshing;
    };


    const getToken = () => {
	    return Token;

    }

    const setToken = (token) => {
        Token = token;
        refreshToken(20);
        return true;
    };

    const ereaseToken = () => {
        Token = null;
        abordRefreshToken();
        window.localStorage.setItem(logoutEventName, Date.now());
        return true;
    }

    return {
        ereaseToken,
        getToken,
        setToken,
	getRefreshedToken,
	waitForTokenRefresh,
    }
};

export default inMemoryJWTManager();
