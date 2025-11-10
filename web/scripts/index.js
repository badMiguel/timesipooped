/**
 * @typedef {{ picture: string, poopTotal: number, failedTotal: number }} UserInfo
 * @typedef {{ poopTotal: number, failedTotal: number }} PoopInfo
 */

/**
 * @param {string} desc
 * @param {string}  [header]
 * */
function showPopupMessage(desc, header) {
    const errorBlur = document.querySelector(".error--blur");
    if (!(errorBlur instanceof HTMLDivElement)) {
        console.error(`failed to find error--blur element.`);
        return;
    }
    const errorContainer = document.querySelector(".error--container");
    if (!(errorContainer instanceof HTMLDivElement)) {
        console.error(`failed to find error--container element.`);
        return;
    }
    const errorDesc = document.querySelector(".error--desc");
    if (!(errorDesc instanceof HTMLParagraphElement)) {
        console.error(`failed to find error--desc element.`);
        return;
    }
    const errorClose = document.querySelector(".error--close");
    if (!(errorClose instanceof HTMLParagraphElement)) {
        console.error(`failed to find error--close element.`);
        return;
    }
    const errorHeader = document.querySelector(".error--header");
    if (!(errorHeader instanceof HTMLHeadingElement)) {
        console.error(`failed to find error--header element.`);
        return;
    }

    if (header) {
        errorHeader.innerText = header;
    }
    errorDesc.innerText = desc;
    errorBlur.style.visibility = "visible";
    errorContainer.style.visibility = "visible";

    errorBlur.addEventListener("click", () => {
        errorBlur.style.visibility = "hidden";
        errorContainer.style.visibility = "hidden";
    });

    errorClose.addEventListener("click", () => {
        errorBlur.style.visibility = "hidden";
        errorContainer.style.visibility = "hidden";
    });
}

/** @returns {Promise<boolean>} */
async function verifyStatus() {
    try {
        const response = await fetch("http://localhost:8001/auth/status", {
            method: "GET",
            credentials: "include",
        });
        if (response.ok) {
            return true;
        } else if (response.status === 401) {
            const getShowLoginPrompt = localStorage.getItem("showLoginPrompt");
            if (!getShowLoginPrompt) {
                localStorage.setItem("showLoginPrompt", "false");
                showPopupMessage(
                    "Please log in to save your poop progress on other devices",
                    "Hi new user!"
                );
                return false;
            }
        } else if (response.status === 403) {
            showPopupMessage("Failed to verify your access.");
            return false;
        }
    } catch (err) {
        console.error(err);
        showPopupMessage("Failed to verify your access.");
    }
    return false;
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

/** @param {number} val */
function updatePoopCounter(val) {
    const poopCounter = document.querySelector(".poop--counter");
    if (!(poopCounter instanceof HTMLHeadingElement)) {
        console.error(`failed to find poop--counter element.`);
        return;
    }
    poopCounter.innerText = val.toString();
}

/**
 * @param {boolean} isPoop
 * @param {boolean} toAdd
 * @returns {PoopInfo | undefined}
 */
function updateToLocalStorage(isPoop, toAdd) {
    const getShowLoginPrompt = localStorage.getItem("showLoginPrompt");
    if (!getShowLoginPrompt) {
        localStorage.setItem("showLoginPrompt", "false");
        showPopupMessage(
            "Please log in to save your poop progress on other devices",
            "Hi new user!"
        );
    }
    let getPoopTotal = localStorage.getItem("poopTotal") || "0";
    let getFailedTotal = localStorage.getItem("failedTotal") || "0";

    if (isPoop) {
        if (toAdd) {
            getPoopTotal = (parseInt(getPoopTotal) + 1).toString();
            localStorage.setItem("poopTotal", getPoopTotal);
        } else {
            if (parseInt(getPoopTotal) < 1) {
                return;
            }
            getPoopTotal = (parseInt(getPoopTotal) - 1).toString();
            localStorage.setItem("poopTotal", getPoopTotal);
        }
    } else {
        if (toAdd) {
            getFailedTotal = (parseInt(getFailedTotal) + 1).toString();
            localStorage.setItem("failedTotal", getFailedTotal);
        } else {
            if (parseInt(getFailedTotal) < 1) {
                return;
            }
            getFailedTotal = (parseInt(getFailedTotal) - 1).toString();
            localStorage.setItem("failedTotal", getFailedTotal);
        }
    }

    /** @type {PoopInfo}*/
    const val = {
        poopTotal: parseInt(getPoopTotal),
        failedTotal: parseInt(getFailedTotal),
    };
    return val;
}

/**
 * @param {boolean} isPoop
 * @param {boolean} toAdd
 * @returns {Promise<PoopInfo | undefined>}
 */
async function updateValueToServer(isPoop, toAdd) {
    // /** @returns {PoopInfo|undefined}*/

    try {
        const isLoggedIn = localStorage.getItem("isLoggedIn");
        if (isLoggedIn) {
            return updateToLocalStorage(isPoop, toAdd);
        }
        let fetchUrl = `http://localhost:8001/${isPoop ? "poop" : "poop/failed"}/${toAdd ? "add" : "sub"}`;
        const response = await fetch(fetchUrl, {
            method: "POST",
            credentials: "include",
        });
        // unauthorized
        if (response.status === 401) {
            return updateToLocalStorage(isPoop, toAdd);
        }
        // forbidden
        if (response.status === 403) {
            showPopupMessage("Failed to verify your access.");
            return;
        }
        // server error
        if (response.status === 500) {
            showPopupMessage("Something went wrong with the server.");
            return;
        }
        return await response.json();
    } catch (err) {
        console.error(
            `Failed ${toAdd ? "add" : "subtract"} <${isPoop ? "poop" : "failed poop"}> value: ${err}`
        );
        showPopupMessage("Failed to update your poop :((");
        return;
    }
}

/** @param {PoopInfo} poopInfo */
async function failedPoop(poopInfo) {
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
        updateValueToServer(false, true);
        poopInfo.failedTotal++;
        updateFailedCounter(poopInfo.failedTotal);
    });
    fPoopSubBtn.addEventListener("click", async () => {
        if (poopInfo.failedTotal > 0) {
            updateValueToServer(false, false);
            poopInfo.failedTotal--;
            updateFailedCounter(poopInfo.failedTotal);
        }
    });
}

/** @param {PoopInfo} poopInfo */
async function poop(poopInfo) {
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
        updateValueToServer(true, true);
        poopInfo.poopTotal++;
        updatePoopCounter(poopInfo.poopTotal);
    });
    poopSubBtn.addEventListener("click", async () => {
        if (poopInfo.poopTotal > 0) {
            updateValueToServer(true, false);
            poopInfo.poopTotal--;
            updatePoopCounter(poopInfo.poopTotal);
        }
    });
}

/** @returns {Promise<UserInfo | undefined>} */
async function fetchVal() {
    try {
        const response = await fetch("http://localhost:8001/get/user", {
            method: "GET",
            credentials: "include",
        });
        if (response.ok) {
            /** @type {UserInfo} */
            const data = await response.json();
            localStorage.setItem("poopTotal", "0");
            localStorage.setItem("failedTotal", "0");
            return data;
        } else if (response.status === 401) {
            const getPoopTotal = localStorage.getItem("poopTotal") || "0";
            const getFailedTotal = localStorage.getItem("failedTotal") || "0";
            if (!getPoopTotal) {
                localStorage.setItem("poopTotal", "0");
            }
            if (!getFailedTotal) {
                localStorage.setItem("failedTotal", "0");
            }
            updatePoopCounter(parseInt(getPoopTotal));
            updateFailedCounter(parseInt(getFailedTotal));
        }
    } catch (err) {
        console.error(err);
        showPopupMessage("Failed to get your poop data :((");
        return undefined;
    }
}

/** @returns {Promise<PoopInfo|undefined>}*/
async function checkStorage() {
    const val = await fetchVal();
    if (val === undefined) {
        return;
    }

    updatePoopCounter(val.poopTotal);
    updateFailedCounter(val.failedTotal);

    const getPicture = localStorage.getItem("picture");
    if (getPicture === null || getPicture != val.picture) {
        localStorage.setItem("picture", val.picture);
    }
    return { poopTotal: val.poopTotal, failedTotal: val.failedTotal };
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

/** @returns {Promise<PoopInfo|undefined>} */
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
    const logoutContainer = document.querySelector(".logout--container");
    if (!(logoutContainer instanceof HTMLDivElement)) {
        console.error(`failed to find logout--container element.`);
        return;
    }
    const logoutClose = document.querySelector(".logout--close");
    if (!(logoutClose instanceof HTMLParagraphElement)) {
        console.error(`failed to find logout--close element.`);
        return;
    }
    const logoutButtonContainer = document.querySelector(".logout--button-container");
    if (!(logoutButtonContainer instanceof HTMLDivElement)) {
        console.error(`failed to find logout--button-container element.`);
        return;
    }
    const mainElement = document.querySelector("body");
    if (!(mainElement instanceof HTMLBodyElement)) {
        console.error(`failed to find logout--button element.`);
        return;
    }
    const errorBlur = document.querySelector(".error--blur");
    if (!(errorBlur instanceof HTMLDivElement)) {
        console.error(`failed to find error--blur element.`);
        return;
    }

    let status = await verifyStatus();
    /** @type {PoopInfo | undefined}*/
    const poopInfo = await checkStorage();
    if (status) {
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
        localStorage.removeItem("isLoggedIn");
    } else {
        profilePicContainer.style.display = "none";
        googleSignIn.style.display = "flex";
        localStorage.setItem("isLoggedIn", "false");
    }

    profileContainer.addEventListener("click", () => {
        if (!status) {
            window.location.href = "http://localhost:8001/auth/login";
        } else {
            if (logoutContainer.style.display === "flex") {
                logoutContainer.style.display = "none";
                errorBlur.style.visibility = "hidden";
                errorBlur.style.zIndex = "2";
            } else {
                logoutContainer.style.display = "flex";
                errorBlur.style.visibility = "visible";
                errorBlur.style.zIndex = "1";
            }
        }
    });

    logoutClose.addEventListener("click", () => {
        logoutContainer.style.display = "none";
        errorBlur.style.visibility = "hidden";
        errorBlur.style.zIndex = "2";
    });

    errorBlur.addEventListener("click", () => {
        logoutContainer.style.display = "none";
        errorBlur.style.visibility = "hidden";
        errorBlur.style.zIndex = "2";
    });

    logoutButtonContainer.addEventListener("click", async () => {
        logoutContainer.style.display = "none";
        errorBlur.style.visibility = "hidden";
        errorBlur.style.zIndex = "2";

        /** @type {Response|undefined} */
        let response = undefined;
        try {
            response = await fetch("http://localhost:8001/auth/logout", {
                credentials: "include",
            });
        } catch (err) {
            showPopupMessage("Server failed to log you out", "Waaaaaaa!");
            console.error(err);
        }

        if (response && response.ok) {
            profilePicContainer.style.display = "none";
            googleSignIn.style.display = "flex";

            localStorage.clear();
            localStorage.setItem("isLoggedIn", "false");
            localStorage.setItem("showLoginPrompt", "false");
            location.reload();
        }

        if (response && !response.ok) {
            showPopupMessage("Server failed to log you out", "Waaaaaaa!");
        }
    });

    return poopInfo;
}

async function main() {
    let poopInfo = await profile();
    if (!poopInfo) {
        const poopTotal = parseInt(localStorage.getItem("poopTotal") || "0");
        const failedTotal = parseInt(localStorage.getItem("failedTotal") || "0");
        poopInfo = { poopTotal: poopTotal, failedTotal: failedTotal };
    }
    poop(poopInfo);
    failedPoop(poopInfo);
}

document.addEventListener("DOMContentLoaded", () => {
    main();
});
