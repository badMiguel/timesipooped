async function addFailPoop() {}

async function addPoop() {}

async function loading() {}

async function fetchVal() {}

async function login() {}

async function main() {
    const profilePic = document.querySelector(".profile-picture--container");
    if (!(profilePic instanceof HTMLDivElement)) {
        console.error(`failed to find profile-picture--container element.`);
        return;
    }
    profilePic.addEventListener("click", async () => {
        await login();
        await fetchVal();
    });

    const poopButton = document.querySelector(".poop--button");
    if (!(poopButton instanceof HTMLParagraphElement)) {
        console.error(`failed to find poop-button element.`);
        return;
    }
    poopButton.addEventListener("click", () => {
        addPoop();
    });

    const fPoopButton = document.querySelector(".failed-poop--button");
    if (!(fPoopButton instanceof HTMLParagraphElement)) {
        console.error(`failed to find failed-poop--button element.`);
        return;
    }
    fPoopButton.addEventListener("click", () => {
        addFailPoop();
    });
}

document.addEventListener("DOMContentLoaded", () => {
    main();
});
