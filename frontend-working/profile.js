import {getUser, getLogoutLink} from './common.js'

function ProfileInit() {
    if (window.location.search.indexOf("?close=true") !== -1) {
      return;
    }

    const loginNav = document.querySelector(".login")
    return getUser().then(user => {
        if (!user) {
            window.location.href = "/index.html";
            return
        }
        var logoutli = document.createElement("li");
        var logoutlink = getLogoutLink()
        logoutli.appendChild(logoutlink);
        loginNav.appendChild(logoutli);
    });
}

ProfileInit().catch(e => {
    console.error(e);
  });