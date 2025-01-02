function addFailPoop() {}

function addPoop() {}

function fetchVal() {}

function main() {
    const poopButton = document.querySelector(".poop--button");
    if (!(poopButton instanceof HTMLParagraphElement)) {
        // todo: add error handling
        console.error(`failed to find poop-button element.`);
        return;
    }

    const fPoopButton = document.querySelector(".failed-poop--button");
    if (!(fPoopButton instanceof HTMLParagraphElement)) {
        // todo: add error handling
        console.error(`failed to find failed-poop--button element.`);
        return;
    }

    poopButton.addEventListener("click", () => {
        addPoop();
    });

    fPoopButton.addEventListener("click", () => {
        addFailPoop();
    });
}

document.addEventListener("DOMContentLoaded", () => {
    fetchVal();
    main();
});
