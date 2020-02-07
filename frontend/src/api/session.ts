
export let session: string;

function getSession() {
    console.log("Getting session header")
    fetch("/session")
    .then(res => res.status == 200 ? res.text() : Promise.reject(res.text())) 
    .then(s => {
        session = s;
    })
    .catch(err => {
        console.log("Could not get session header:", err)
        setTimeout(getSession, 10000);
    })
}

getSession();