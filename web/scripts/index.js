function fetchError() {}

function notAuthenticated() {
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
            notAuthenticated();
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

async function fetchVal() {}

async function profile() {
    try {
        await verifyStatus();
    } catch (err) {}

    const profileContainer = document.querySelector(".profile--container");
    if (!(profileContainer instanceof HTMLDivElement)) {
        console.error(`failed to find profile--container element.`);
        return;
    }
    profileContainer.addEventListener("click", () => {
        window.location.href = "http://localhost:8081/auth/login";
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
