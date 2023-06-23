import {getLogoutLink, getUser} from './common.js'

function main() {
    if (window.location.search.indexOf("?close=true") !== -1) {
      return;
    }
    const loginNav = document.querySelector(".login")
    const navlist = document.querySelector(".login").getElementsByTagName("li");
    return getUser().then(user => {
        if (!user) {
            return
        }
        var username = user.name
        navlist[0].innerHTML = `<a href="profile.html">${username} </a>`;
        var logoutli = document.createElement("li");
        var logoutlink = getLogoutLink()
        logoutli.appendChild(logoutlink);
        loginNav.appendChild(logoutli);
    });
}

main().catch(e => {
    console.error(e);
  });