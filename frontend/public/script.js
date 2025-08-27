window.addEventListener("load", function () {
    // Header Menu toggles
    const showMenuBtn = document.querySelector("#showMenu");
    const hideMenuBtn = document.querySelector("#hideMenu");
    const mobileNav = document.querySelector("#mobileNav"); 
    showMenuBtn.addEventListener("click", function () {
    // reset classes before showing
        mobileNav.classList.remove("hidden", "animate-fade-out");
        mobileNav.classList.add("animate-fade-in");
    }); 
    hideMenuBtn.addEventListener("click", function () {
        // trigger fade-out
        mobileNav.classList.remove("animate-fade-in");
        mobileNav.classList.add("animate-fade-out");  
        // after fade-out finishes, hide completely
        mobileNav.addEventListener(
        "animationend",
        () => {
            if (mobileNav.classList.contains("animate-fade-out")) {
            mobileNav.classList.add("hidden");
            }
        },
        { once: true }
      );
    }); 
    // FAQ section
    document.querySelectorAll("[toggleElement]").forEach((toggle) => {
      toggle.addEventListener("click", function () {
        const answerElement = toggle.querySelector("[answer]");
        const caretElement = toggle.querySelector("img");   
        if (answerElement.classList.contains("hidden")) {
            answerElement.classList.remove("hidden");
            caretElement.classList.add("rotate-90");
        } else {
            answerElement.classList.add("hidden");
            caretElement.classList.remove("rotate-90");
        }
      });
    });
    // Loading screen
    const loader = document.getElementById("loadingScreen");
    const app = document.getElementById("app");

    // Wait 3.5 seconds, then fade out loader and show app
    setTimeout(() => {
        loader.classList.add("animate-fade-out");

        loader.addEventListener("animationend", () => {
        loader.style.display = "none";
        app.classList.remove("hidden");
        app.classList.add("animate-fade-in");
        }, { once: true });
    }, 2800);
});
