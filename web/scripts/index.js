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
    // TODO
    console.log("not authenticated");
}

function updateValueError() {}

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

/**
 * @param {boolean} isPoop
 * @param {boolean} toAdd
 */
async function updateValue(isPoop, toAdd) {
    try {
        let fetchUrl = `http://localhost:8081/${isPoop ? "poop" : "poop/failed"}/${toAdd ? "add" : "sub"}`;
        const response = await fetch(fetchUrl, {
            method: "POST",
            credentials: "include",
        });
    } catch (err) {
        console.error(
            `Failed ${toAdd ? "add" : "subtract"} <${isPoop ? "poop" : "failed poop"}> value: ${err}`
        );
        updateValueError();
        return;
    }
}

async function failedPoop() {
    const fPoopAddBtn = document.querySelector(".failed-poop-add--button");
    if (!(fPoopAddBtn instanceof HTMLParagraphElement)) {
        console.error(`failed to find failed-poop-add--button element.`);
        return;
    }
    const fPoopSubBtn = document.querySelector(".failed-poop-sub--button");
    if (!(fPoopSubBtn instanceof HTMLParagraphElement)) {
        console.error(`failed to find failed-poop-sub--button element.`);
        return;
    }
    fPoopAddBtn.addEventListener("click", async () => {
        await updateValue(false, true);
    });
    fPoopSubBtn.addEventListener("click", async () => {
        await updateValue(false, false);
    });
}

async function poop() {
    const poopAddBtn = document.querySelector(".poop-add--button");
    if (!(poopAddBtn instanceof HTMLParagraphElement)) {
        console.error(`failed to find poop-add--button element.`);
        return;
    }
    const poopSubBtn = document.querySelector(".poop-sub--button");
    if (!(poopSubBtn instanceof HTMLParagraphElement)) {
        console.error(`failed to find poop-sub--button element.`);
        return;
    }

    poopAddBtn.addEventListener("click", async () => {
        await updateValue(true, true);
    });
    poopSubBtn.addEventListener("click", async () => {
        await updateValue(true, false);
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

/** @returns {boolean} */
function isPictureExpired() {
    const photoExpireDate = localStorage.getItem("picture_timestamp");
    if (photoExpireDate === null) {
        const now = new Date().setHours(0, 0, 0, 0).toString();
        localStorage.setItem("picture_timestamp", now);
        return true;
    }
    const past = new Date(parseInt(photoExpireDate)).getTime();
    const now = new Date().setHours(0, 0, 0, 0);
    if (past < now) {
        localStorage.setItem("picture_timestamp", now.toString());
        return true;
    }
    return false;
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
            if (isPictureExpired()) {
                image.src = cacheBustedSrc;
                image.onload = () => {
                    profilePic.src = cacheBustedSrc;
                };
            } else {
                profilePic.src = getPic;
            }
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
