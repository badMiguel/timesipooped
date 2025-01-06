/**
 * @typedef {{
 * picture: string,
 * poopTotal: number,
 * failedTotal: number,
 * }} UserInfo
 *
 * @typedef {{
 * poopTotal: number
 * failedTotal: number
 * }} UpdatedPoop
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
 * @returns {Promise<UpdatedPoop | undefined>}
 */
async function updateValue(isPoop, toAdd) {
    try {
        let fetchUrl = `http://localhost:8081/${isPoop ? "poop" : "poop/failed"}/${toAdd ? "add" : "sub"}`;
        const response = await fetch(fetchUrl, {
            method: "POST",
            credentials: "include",
        });
        if (!response.ok) {
            throw new Error(`Response not ok. Status code: ${response.status}`);
        }
        return await response.json();
    } catch (err) {
        console.error(
            `Failed ${toAdd ? "add" : "subtract"} <${isPoop ? "poop" : "failed poop"}> value: ${err}`
        );
        updateValueError();
        return;
    }
}

/** @param {number} val */
function updateFailedCounter(val) {
    const failedCounter = document.querySelector(".failed-poop--counter");
    if (!(failedCounter instanceof HTMLHeadingElement)) {
        console.error(`failed to find failed-poop--counter element.`);
        return;
    }
    failedCounter.innerText = val.toString();
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
        const val = await updateValue(false, true);
        if (val) {
            updateFailedCounter(val.failedTotal);
        }
    });
    fPoopSubBtn.addEventListener("click", async () => {
        const val = await updateValue(false, false);
        if (val) {
            updateFailedCounter(val.failedTotal);
        }
    });
}

/** @param {number} val */
function updatePoopCounter(val) {
    const poopCounter = document.querySelector(".poop--counter");
    if (!(poopCounter instanceof HTMLHeadingElement)) {
        console.error(`failed to find poop--counter element.`);
        return;
    }
    poopCounter.innerText = val.toString();
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
        const val = await updateValue(true, true);
        if (val) {
            updatePoopCounter(val.poopTotal);
        }
    });
    poopSubBtn.addEventListener("click", async () => {
        const val = await updateValue(true, false);
        if (val) {
            updatePoopCounter(val.poopTotal);
        }
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
        console.error(err);
        fetchError();
        return null;
    }
}

async function checkStorage() {
    const val = await fetchVal();
    if (val === null) {
        return;
    }

    updatePoopCounter(val.poopTotal);
    updateFailedCounter(val.failedTotal);

    const getPicture = localStorage.getItem("picture");
    if (getPicture === null) {
        localStorage.setItem("picture", val.picture);
        return;
    }
    if (getPicture !== val.picture) {
        localStorage.setItem("picture", val.picture);
        return;
    }
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
                    profilePic.src = image.src;
                    localStorage.setItem("picture", image.src);
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
