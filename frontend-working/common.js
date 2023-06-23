function getCookies() {
    console.log("doc cookier --",document.cookie.split("; "))

    return document.cookie.split("; ").reduce((c, x) => {
      const splitted = x.split("=");
      c[splitted[0]] = splitted[1];
      return c;
    }, {});
  }
  

function req(endpoint, data = {}) {
    const cloneData = Object.assign({}, data);
    const cookies = getCookies();
    const token = cookies["XSRF-TOKEN"];
    console.log("cookies --> ", cookies)

    if (cloneData.hasOwnProperty("headers")) {
      cloneData.headers = new Headers(cloneData.headers);
    } else {
      cloneData.headers = new Headers();
    }
  
    if (token) {
      cloneData.headers.append("X-XSRF-TOKEN", token);
    }
    
    return fetch(endpoint, cloneData).then(resp => {
      if (resp.status >= 400) {
        throw resp;
      }
      console.log("requsted ", endpoint)
      return resp.json().catch(() => null);
    });
  }


  function getUser() {
    return req("/auth/user").catch(e => {
        if (e.status && e.status === 401) return null;
        throw e;
    });
}



function getProviders() {
    return req("/auth/list");
}

function getLogoutLink() {
    const a = document.createElement("a");
    a.href = "#";
    a.textContent = "Logout";
    a.className = "login__prov";
    a.addEventListener("click", e => {
      e.preventDefault();
      req("/auth/logout")
        .then(() => {
          window.location.replace(window.location.href);
        })
        .catch(errorHandler);
    });
    return a;
}

function errorHandler(err) {
    // const status = document.querySelector(".status__label");
    if (err instanceof Response) {
      err.text().then(text => {
        try {
          const data = JSON.parse(text);
          if (data.error) {
            // status.textContent = data.error;
            console.error(data.error);
            return;
          }
        } catch {
        }
        // status.textContent = text;
        console.error(text);
      });
      return;
    }
    // status.textContent = err.message;
    console.error(err.message);
  }


export {getUser}
export {getProviders}
export {getLogoutLink}
export {req}
export {errorHandler}

