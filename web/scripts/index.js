/**
 * @typedef {{
 * family_name: string,
 * given_name: string,
 * picture: string,
 * poop_total: number,
 * failed_total: number,
 * }} UserInfo
 */

function fetchError() {}

function authError() {
    console.log("not authenticated");
}

/**
 * @returns {Promise<boolean>}
 */
async function verifyStatus() {
    try {
        const response = await fetch("http://localhost:8081/auth/status", {
            method: "GET",
            credentials: "include",
        });
        if (!response.ok) {
            authError();
            return false;
        }
    } catch (err) {
        console.log(err);
        fetchError();
        return false;
    }
    return true;
}

async function failedPoop() {
    const fPoopButton = document.querySelector(".failed-poop--button");
    if (!(fPoopButton instanceof HTMLParagraphElement)) {
        console.error(`failed to find failed-poop--button element.`);
        return;
    }
    fPoopButton.addEventListener("click", () => {
        try {
        } catch (err) {}
    });
}

async function poop() {
    const poopButton = document.querySelector(".poop--button");
    if (!(poopButton instanceof HTMLParagraphElement)) {
        console.error(`failed to find poop-button element.`);
        return;
    }
    poopButton.addEventListener("click", async () => {
        await verifyStatus();
    });
}

async function loading() {}

/**
 * @returns {Promise<UserInfo | null>}
 */
async function fetchVal() {
    try {
        const response = await fetch("http://localhost:8081/get/user", {
            method: "GET",
            credentials: "include",
        });
        /** @type {UserInfo} */
        const data = await response.json();
        return data;
    } catch (err) {
        console.log(err);
        fetchError();
        return null;
    }
}

/**
 * @param  {string} key
 * @param  {string | number} item
 */
function checkStorageHelper(key, item) {
    const getInfo = localStorage.getItem(key);
    if (getInfo === null) {
        if (typeof item === "number") {
            localStorage.setItem(key, item.toString());
        } else {
            localStorage.setItem(key, item);
        }
        return;
    }
    if (getInfo !== item) {
        if (typeof item === "number") {
            localStorage.setItem(key, item.toString());
        } else {
            localStorage.setItem(key, item);
        }
        return;
    }
}

async function checkStorage() {
    const val = await fetchVal();
    if (val === null) {
        return;
    }

    checkStorageHelper("given_name", val.given_name);
    checkStorageHelper("family_name", val.family_name);
    checkStorageHelper("picture", val.picture);
    checkStorageHelper("poop_total", val.poop_total);
    checkStorageHelper("failed_total", val.failed_total);
}

async function profile() {
    const profileContainer = document.querySelector(".profile--container");
    if (!(profileContainer instanceof HTMLDivElement)) {
        console.error(`failed to find profile--container element.`);
        return;
    }
    const profilePicContainer = document.querySelector(".profile-picture--container");
    if (!(profilePicContainer instanceof HTMLDivElement)) {
        console.error(`failed to find profile-pic--container element.`);
        return;
    }
    const profilePic = document.querySelector(".profile-picture");
    if (!(profilePic instanceof HTMLImageElement)) {
        console.error(`failed to find .profile-picture element.`);
        return;
    }
    const googleSignIn = document.querySelector(".google-sign-in");
    if (!(googleSignIn instanceof HTMLDivElement)) {
        console.error(`failed to find google-sign-in element.`);
        return;
    }

    const status = await verifyStatus();
    if (status) {
        await checkStorage();
        const getPic = localStorage.getItem("picture");
        if (getPic !== null) {
            const image = new Image();
            const cacheBustedSrc = `${getPic}?t=${new Date().getTime()}`;
            image.src = cacheBustedSrc;
            image.onload = () => {
                profilePic.src = cacheBustedSrc;
            };
        }
        profilePicContainer.style.display = "flex";
        googleSignIn.style.display = "none";
    } else {
        profilePicContainer.style.display = "none";
        googleSignIn.style.display = "flex";
    }

    profileContainer.addEventListener("click", () => {
        if (!status) {
            window.location.href = "http://localhost:8081/auth/login";
        } else {
            // TODO add sign out
        }
    });
}

async function main() {
    profile();
    poop();
    failedPoop();
}

document.addEventListener("DOMContentLoaded", () => {
    main();
});
