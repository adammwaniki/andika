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
            // Initialize auto-resize after notes are visible
            setTimeout(() => {
              initAutoResize();
            }, 100);
          }
        },
        { once: true }
      );
    }, 2800);
  } else {
    // If no loading screen, initialize immediately
    setTimeout(() => {
      initAutoResize();
    }, 100);
  }

  // Improved auto-resize textarea function
  function autoResizeTextarea(textarea) {
    // Store the original height
    const originalHeight = textarea.style.height;
    
    // Reset height to auto to get the actual scrollHeight
    textarea.style.height = 'auto';
    
    // Calculate the new height based on content
    let newHeight = textarea.scrollHeight;
    
    // For single-row textareas (like titles), ensure minimum height
    if (textarea.getAttribute('rows') === '1') {
      const lineHeight = parseInt(getComputedStyle(textarea).lineHeight) || 24;
      newHeight = Math.max(newHeight, lineHeight);
    }
    
    // Set the new height
    textarea.style.height = newHeight + 'px';
    
    // If height changed significantly, trigger a reflow
    if (Math.abs(newHeight - parseInt(originalHeight || '0')) > 5) {
      textarea.offsetHeight; // Force reflow
    }
  }

  // Validation helper functions
  function showTitleValidationError(titleTextarea) {
    // Check if error message already exists
    const existingError = titleTextarea.parentElement.querySelector('.title-error-message');
    if (existingError) {
      return; // Error already showing
    }
    
    // Create error message element
    const errorMessage = document.createElement('div');
    errorMessage.className = 'title-error-message';
    errorMessage.textContent = 'A title is necessary in order to save';
    
    // Insert error message after the textarea
    titleTextarea.parentElement.insertBefore(errorMessage, titleTextarea.nextSibling);
    
    // Add error state class to textarea
    titleTextarea.classList.add('error-state');
    
    // Focus on the title textarea
    titleTextarea.focus();
  }
  
  function removeTitleValidationError(titleTextarea) {
    // Remove error message if it exists
    const errorMessage = titleTextarea.parentElement.querySelector('.title-error-message');
    if (errorMessage) {
      errorMessage.remove();
    }
    
    // Remove error state class
    titleTextarea.classList.remove('error-state');
  }

  // Delete note card function
  function deleteNoteCard(noteElement) {
    // Close expanded view if open
    if (noteElement.classList.contains("expanded") || noteElement.classList.contains("expanding")) {
      closeExpandedNote(noteElement);
    }
    
    // Add deleting class for animation
    noteElement.classList.add('deleting');
    
    // Remove from DOM after animation completes
    setTimeout(() => {
      noteElement.remove();
      
      // Here we make an API call to delete the note
      // Extract the note ID from a data attribute or the content
      // i.e,:
      // const noteId = noteElement.getAttribute('data-note-id');
      // fetch(`/api/notes/${noteId}`, { method: 'DELETE' })
      //   .catch(error => {
      //     console.error('Error deleting note:', error);
      //     // Optionally show error message to user
      //   });
      
      console.log('Note deleted from UI - API call will be made here');
      
    }, 400); // Match the CSS transition duration
  }
  
  // Clear create note form function
  function clearCreateNoteForm(createNoteElement) {
    // Get the textareas
    const titleTextarea = createNoteElement.querySelector('#inputTitle');
    const contentTextarea = createNoteElement.querySelector('#inputContent');
    
    // Clear the values
    if (titleTextarea) {
      titleTextarea.value = '';
      removeTitleValidationError(titleTextarea);
      autoResizeTextarea(titleTextarea);
    }
    
    if (contentTextarea) {
      contentTextarea.value = '';
      autoResizeTextarea(contentTextarea);
    }
    
    // Close expanded view if open
    if (createNoteElement.classList.contains("expanded") || createNoteElement.classList.contains("expanding")) {
      closeExpandedNote(createNoteElement);
    }
    
    // Focus on title for better UX
    if (titleTextarea) {
      titleTextarea.focus();
    }
  }

  // Apply auto-resize to all textareas
  function initAutoResize() {
    document.querySelectorAll('textarea').forEach(textarea => {
      // Initial resize for textareas with content
      if (textarea.value.trim() !== '' || textarea.textContent.trim() !== '') {
        // Small delay to ensure proper rendering
        setTimeout(() => {
          autoResizeTextarea(textarea);
        }, 10);
      } else {
        // For empty textareas, still set minimum height
        autoResizeTextarea(textarea);
      }
      
      // Add input listener
      textarea.addEventListener('input', () => autoResizeTextarea(textarea));
      
      // Add focus listener to ensure proper height
      textarea.addEventListener('focus', () => {
        setTimeout(() => autoResizeTextarea(textarea), 10);
      });
    });
  }

  // Note expansion
  const notes = document.querySelectorAll(".note-card");
  const section = document.querySelector(".sectionSizeNotes");

  notes.forEach(note => {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const trashBtn = note.querySelector(".trashBtn");
    //Ellipsis button will be used to open a menu overlay above the card it is on
    const ellipsisBtn = note.querySelector(".fa-ellipsis")?.parentElement;
    const saveBtn = note.querySelector(".saveNote");
    const saveNewBtn = note.querySelector("#saveNewNote, [id^='save']");

    // Hide close button by default
    if (closeBtn && !closeBtn.classList.contains("hidden")) {
      closeBtn.classList.add("hidden");
    }

    // Initialize textareas in this note
    const textareas = note.querySelectorAll('textarea');
    textareas.forEach(textarea => {
      // Ensure proper initial height
      setTimeout(() => {
        autoResizeTextarea(textarea);
      }, 50);
    });

    // Add real-time validation for title textarea in CreateNote
    const titleTextarea = note.querySelector('#inputTitle');
    if (titleTextarea) {
      titleTextarea.addEventListener('input', () => {
        const titleValue = titleTextarea.value.trim();
        if (titleValue !== '') {
          // Remove validation error if user starts typing
          removeTitleValidationError(titleTextarea);
        }
      });
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
        
        // For CreateNote card, just clear the form instead of deleting
        if (note.querySelector('#inputTitle')) {
          clearCreateNoteForm(note);
          return;
        }
        
        // For existing notes, delete with animation
        deleteNoteCard(note);
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
        
        // Get the title textarea for this note
        const titleTextarea = note.querySelector('#inputTitle');
        if (titleTextarea) {
          const titleValue = titleTextarea.value.trim();
          
          if (titleValue === '') {
            // Show validation error
            showTitleValidationError(titleTextarea);
            return; // Don't proceed with save
          } else {
            // Remove any existing validation error
            removeTitleValidationError(titleTextarea);
          }
        }
        
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
        // Re-apply auto-resize
        setTimeout(() => autoResizeTextarea(titleTextarea), 10);
      }
    }

    if (contentArea) {
      contentArea.classList.remove('content-area');
      const contentTextarea = contentArea.querySelector('textarea');
      if (contentTextarea) {
        contentTextarea.classList.remove('content-textarea', 'auto-resize');
        contentTextarea.setAttribute('rows', '3');
        contentTextarea.style.height = 'auto';
        // Re-apply auto-resize
        setTimeout(() => autoResizeTextarea(contentTextarea), 10);
      }
    }

    if (buttonArea) {
      buttonArea.classList.remove('button-area');
    }
  }

  // Function to expand note smoothly
  function expandNote(note) {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis")?.parentElement;

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
      
      // Re-initialize auto-resize for textareas in expanded note
      const textareas = note.querySelectorAll('textarea');
      textareas.forEach(textarea => {
        setTimeout(() => autoResizeTextarea(textarea), 50);
      });
    }, 500);
  }

  // Function to close expanded note instantly
  function closeExpandedNote(note) {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis")?.parentElement;

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

  // Fallback initialization if no loading screen
  if (!loader) {
    initAutoResize();
  }
});