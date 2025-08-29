window.addEventListener("load", function () {
    // Header Menu toggles
    const showMenuBtn = document.querySelector("#showMenu");
    const hideMenuBtn = document.querySelector("#hideMenu");
    const mobileNav = document.querySelector("#mobileNav");
     if (showMenuBtn && hideMenuBtn && mobileNav) { 
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
     }
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
        const appNotes = document.getElementById("appNotes");

        if (loader) {
    setTimeout(() => {
      loader.classList.add("animate-fade-out");

      loader.addEventListener(
        "animationend",
        () => {
          loader.style.display = "none";

          if (app) {
            app.classList.remove("hidden");
            app.classList.add("animate-fade-in");
          } else if (appNotes) {
            appNotes.classList.remove("hidden");
            appNotes.classList.add("animate-fade-in");
          }
        },
        { once: true }
      );
    }, 2800);
  }

  // Auto-resize textarea function
  function autoResizeTextarea(textarea) {
    textarea.style.height = 'auto';
    textarea.style.height = textarea.scrollHeight + 'px';
  }

  // Apply auto-resize to all textareas
  function initAutoResize() {
    document.querySelectorAll('textarea').forEach(textarea => {
      autoResizeTextarea(textarea);
      textarea.addEventListener('input', () => autoResizeTextarea(textarea));
    });
  }

  // Note expansion
  const notes = document.querySelectorAll(".note-card");
  const section = document.querySelector(".sectionSizeNotes");

  notes.forEach(note => {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const trashBtn = note.querySelector(".trashBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis-vertical")?.parentElement;
    const saveBtn = note.querySelector(".saveNote");
    const saveNewBtn = note.querySelector("#saveNewNote, [id^='save']");

    // Hide close button by default
    if (closeBtn && !closeBtn.classList.contains("hidden")) {
      closeBtn.classList.add("hidden");
    }

    // Expand note on click
    note.addEventListener("click", (e) => {
      // Prevent clicks on buttons or icons from triggering expansion
      if (
        e.target.closest(".closeNoteBtn") ||
        e.target.closest(".trashBtn") ||
        e.target.closest("button") ||
        e.target.tagName.toLowerCase() === "button" ||
        e.target.closest("saveNote")
      ) return;

      // Close all other notes first
      notes.forEach(n => {
        if (n !== note && n.classList.contains("expanded")) {
          closeExpandedNote(n);
        }
      });

      // Expand this note smoothly
      expandNote(note);
    });

    // Close button click handler
    if (closeBtn) {
      closeBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        closeExpandedNote(note);
      });
    }

    // Trash button click handler
    if (trashBtn) {
      trashBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        // Delete logic can be called here from the API
        // For now, we just close the expanded view if it's open
        if (note.classList.contains("expanded") || note.classList.contains("expanding")) {
          closeExpandedNote(note);
        }
      });
    }

    // Save existing note button click handler
    if (saveBtn) {
      saveBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        // Saving logic can be called here from the API
        // For now, we just close the expanded view if it's open
        if (note.classList.contains("expanded") || note.classList.contains("expanding")) {
          closeExpandedNote(note);
        }
      });
    }

    // Save new note button click handler
    if (saveNewBtn) {
      saveNewBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        // Saving logic can be called here from the API
        // For now, we just close the expanded view if it's open
        if (note.classList.contains("expanded") || note.classList.contains("expanding")) {
          closeExpandedNote(note);
        }
      });
    }
  });

  // Function to restructure note for expanded view
  function restructureExpandedNote(note) {
    const titleArea = note.querySelector('textarea[rows="1"]')?.parentElement;
    const contentArea = note.querySelector('textarea[rows="3"]')?.parentElement;
    const buttonArea = note.querySelector('.flex.justify-between');

    if (titleArea) {
      titleArea.classList.add('title-area');
      const titleTextarea = titleArea.querySelector('textarea');
      if (titleTextarea) {
        titleTextarea.classList.add('title-textarea');
      }
    }

    if (contentArea) {
      contentArea.classList.add('content-area');
      const contentTextarea = contentArea.querySelector('textarea');
      if (contentTextarea) {
        contentTextarea.classList.add('content-textarea', 'auto-resize');
        // Calculate minimum rows based on viewport height
        const minRows = Math.floor((window.innerHeight - 300) / 24); // Approximate line height
        contentTextarea.setAttribute('rows', Math.max(minRows, 10));
      }
    }

    if (buttonArea) {
      buttonArea.classList.add('button-area');
    }
  }

  // Function to reset note structure
  function resetNoteStructure(note) {
    const titleArea = note.querySelector('.title-area');
    const contentArea = note.querySelector('.content-area');
    const buttonArea = note.querySelector('.button-area');

    if (titleArea) {
      titleArea.classList.remove('title-area');
      const titleTextarea = titleArea.querySelector('textarea');
      if (titleTextarea) {
        titleTextarea.classList.remove('title-textarea');
        titleTextarea.style.height = 'auto';
      }
    }

    if (contentArea) {
      contentArea.classList.remove('content-area');
      const contentTextarea = contentArea.querySelector('textarea');
      if (contentTextarea) {
        contentTextarea.classList.remove('content-textarea', 'auto-resize');
        contentTextarea.setAttribute('rows', '3');
        contentTextarea.style.height = 'auto';
      }
    }

    if (buttonArea) {
      buttonArea.classList.remove('button-area');
    }
  }

  // Function to expand note smoothly
  function expandNote(note) {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis-vertical")?.parentElement;

    // Get original dimensions and position
    const rect = note.getBoundingClientRect();
    const originalStyles = {
      width: rect.width + 'px',
      height: rect.height + 'px',
      top: rect.top + 'px',
      left: rect.left + 'px',
      padding: getComputedStyle(note).padding,
      margin: getComputedStyle(note).margin,
      borderRadius: getComputedStyle(note).borderRadius
    };

    // Add expanding class and set initial position
    note.classList.add('expanding');
    Object.assign(note.style, originalStyles);

    // Show close button, hide ellipsis immediately
    if (closeBtn) closeBtn.classList.remove("hidden");
    if (ellipsisBtn) ellipsisBtn.classList.add("hidden");

    // Force a reflow
    note.offsetHeight;

    // Animate to fullscreen
    requestAnimationFrame(() => {
      note.style.top = '0px';
      note.style.left = '0px';
      note.style.width = '100vw';
      note.style.height = '100vh';
      note.style.margin = '0';
      note.style.padding = '2rem';
      note.style.borderRadius = '0';
    });

    // After transition completes, add expanded class and restructure
    setTimeout(() => {
      note.classList.remove('expanding');
      note.classList.add('expanded');
      
      // Clear inline styles since expanded class handles them
      note.style.cssText = '';
      
      // Restructure the note for expanded view
      restructureExpandedNote(note);
      
      // Initialize auto-resize for textareas in expanded note
      initAutoResize();
    }, 500);
  }

  // Function to close expanded note instantly
  function closeExpandedNote(note) {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis-vertical")?.parentElement;

    // If not expanded, just return
    if (!note.classList.contains("expanded") && !note.classList.contains("expanding")) {
      return;
    }

    // Reset structure first
    resetNoteStructure(note);

    // Remove all classes and clear styles immediately
    note.classList.remove('expanded', 'expanding');
    note.style.cssText = ''; // Clear all inline styles
    
    // Reset button visibility
    if (closeBtn) closeBtn.classList.add("hidden");
    if (ellipsisBtn) ellipsisBtn.classList.remove("hidden");
  }

  // ESC key closes any expanded note
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
      notes.forEach(note => {
        if (note.classList.contains("expanded") || note.classList.contains("expanding")) {
          closeExpandedNote(note);
        }
      });
    }
  });

  // Handle window resize for expanded notes
  window.addEventListener("resize", () => {
    notes.forEach(note => {
      if (note.classList.contains("expanded")) {
        const contentTextarea = note.querySelector('.content-textarea');
        if (contentTextarea) {
          const minRows = Math.floor((window.innerHeight - 300) / 24);
          contentTextarea.setAttribute('rows', Math.max(minRows, 10));
          autoResizeTextarea(contentTextarea);
        }
      } else if (note.classList.contains("expanding")) {
        // If note is expanding during resize, complete the transition immediately
        note.style.width = '100vw';
        note.style.height = '100vh';
      }
    });
  });

  // Initialize auto-resize for existing textareas
  initAutoResize();
});