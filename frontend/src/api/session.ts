import { API_URL } from "./url";

export let session: {
    Session: string
};

function getSession() {
    console.log("Getting session header")
    fetch(`${API_URL}/app/session`)
    .then(res => res.status == 200 ? res.text() : Promise.reject(res.text())) 
    .then(s => {
        session = {
            Session: s,
        };
    })
    .catch(err => {
        console.log("Could not get session header:", err)
        setTimeout(getSession, 10000);
    })
}

getSession();