import {getUser, errorHandler, getProviders} from './common.js'


function login(prov) {
    console.log("logging in --", prov)
    return new Promise((resolve, reject) => {
      const url = window.location.href + "?close=true";
      const eurl = encodeURIComponent(url);
      const win = window.open(
        "/auth/" + prov + "/login?id=shellshare-dev&from=" + eurl
      );
      const interval = setInterval(() => {
        try {
          if (win.closed) {
            reject(new Error("Login aborted"));
            clearInterval(interval);
            return;
          }
          if (win.location.search.indexOf("error") !== -1) {
            reject(new Error(win.location.search));
            win.close();
            clearInterval(interval);
            return;
          }
          if (win.location.href.indexOf(url) === 0) {
            resolve();
            win.close();
            clearInterval(interval);
          }
        } catch (e) {
        }
      }, 100);
    });
  }

function LoginInit() {
    if (window.location.search.indexOf("?close=true") !== -1) {
    //   document.body.textContent = "Logged in!";
      return;
    }
    console.log("herhere")
    var logincontainer = document.querySelector(".container")
    // console.log(logincontainer)

    var loginbuttons = logincontainer.getElementsByClassName("login-btn")
    console.log(loginbuttons)

    // navlist = document.querySelector(".login").getElementsByTagName("li");
    return getUser().then(user => {
        if (user) {
            console.log("logged in user")
            // win.open("/profile.html")
            window.location.href = "/profile.html";
            return
        }
        console.log("user ->", user)
        let formSwitcher = () => {
        };
        getProviders().then(providers =>
            providers.map(prov => {
                console.log(prov)
            }))
        for (let i=0 ;i<loginbuttons.length; i++){
            let anchor = loginbuttons[i].getElementsByTagName('a')
            let service = anchor[0].id
            anchor[0].addEventListener("click", e => {
                formSwitcher();
                e.preventDefault();
                login(service)
                  .then(() => {
                    window.location.replace(window.location.href);
                  })
                  .catch(errorHandler);
              });
        }
    });
}

LoginInit().catch(e => {
    console.error(e);
  });