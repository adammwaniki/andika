// Load API integration first
window.addEventListener("load", function () {
    // API Integration functions for Andika frontend
    const API_BASE_URL = '/api/v1';

    const apiCall = async (endpoint, options = {}) => {
      const url = `${API_BASE_URL}${endpoint}`;
      const config = {
        headers: {
          'Content-Type': 'application/json',
        },
        ...options,
      };
      
      try {
        const response = await fetch(url, config);
        const data = await response.json();
        
        if (!response.ok) {
          throw new Error(data.error || `HTTP error! status: ${response.status}`);
        }
        
        return data;
      } catch (error) {
        console.error('API call failed:', error);
        throw error;
      }
    };

    const noteAPI = {
      create: async (name, content) => {
        return await apiCall('/notes', {
          method: 'POST',
          body: JSON.stringify({ name, content }),
        });
      },
      list: async () => {
        return await apiCall('/notes');
      },
      get: async (id) => {
        return await apiCall(`/notes/${id}`);
      },
      update: async (id, content, mode = 'overwrite') => {
        return await apiCall(`/notes/${id}`, {
          method: 'PUT',
          body: JSON.stringify({ mode, content }),
        });
      },
      delete: async (id) => {
        return await apiCall(`/notes/${id}`, {
          method: 'DELETE',
        });
      },
    };

    const showNotification = (message, type = 'info') => {
      const notification = document.createElement('div');
      notification.className = `notification ${type}`;
      notification.textContent = message;
      notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 12px 20px;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        z-index: 1000;
        transition: opacity 0.3s ease;
        ${type === 'error' ? 'background-color: #DC2626;' : 
          type === 'success' ? 'background-color: #16A34A;' : 
          'background-color: #2563EB;'}
      `;
      
      document.body.appendChild(notification);
      
      setTimeout(() => {
        notification.style.opacity = '0';
        setTimeout(() => notification.remove(), 300);
      }, 3000);
    };

    // Header Menu toggles
    const showMenuBtn = document.querySelector("#showMenu");
    const hideMenuBtn = document.querySelector("#hideMenu");
    const mobileNav = document.querySelector("#mobileNav");
     if (showMenuBtn && hideMenuBtn && mobileNav) { 
        showMenuBtn.addEventListener("click", function () {
            mobileNav.classList.remove("hidden", "animate-fade-out");
            mobileNav.classList.add("animate-fade-in");
        }); 
        hideMenuBtn.addEventListener("click", function () {
            mobileNav.classList.remove("animate-fade-in");
            mobileNav.classList.add("animate-fade-out");  
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
                setTimeout(() => {
                  initAutoResize();
                }, 100);
              }
            },
            { once: true }
          );
        }, 2800);
    } else {
        setTimeout(() => {
          initAutoResize();
        }, 100);
    }

    // Auto-resize textarea function
    function autoResizeTextarea(textarea) {
        const originalHeight = textarea.style.height;
        textarea.style.height = 'auto';
        let newHeight = textarea.scrollHeight;
        
        if (textarea.getAttribute('rows') === '1') {
          const lineHeight = parseInt(getComputedStyle(textarea).lineHeight) || 24;
          newHeight = Math.max(newHeight, lineHeight);
        }
        
        textarea.style.height = newHeight + 'px';
        
        if (Math.abs(newHeight - parseInt(originalHeight || '0')) > 5) {
          textarea.offsetHeight;
        }
    }

    // Validation helper functions
    function showTitleValidationError(titleTextarea) {
        const existingError = titleTextarea.parentElement.querySelector('.title-error-message');
        if (existingError) {
          return;
        }
        
        const errorMessage = document.createElement('div');
        errorMessage.className = 'title-error-message';
        errorMessage.textContent = 'A title is necessary in order to save';
        
        titleTextarea.parentElement.insertBefore(errorMessage, titleTextarea.nextSibling);
        titleTextarea.classList.add('error-state');
        titleTextarea.focus();
    }
    
    function removeTitleValidationError(titleTextarea) {
        const errorMessage = titleTextarea.parentElement.querySelector('.title-error-message');
        if (errorMessage) {
          errorMessage.remove();
        }
        titleTextarea.classList.remove('error-state');
    }

    // API-integrated functions
    async function handleSaveNewNote(noteElement) {
        const titleTextarea = noteElement.querySelector('#inputTitle');
        const contentTextarea = noteElement.querySelector('#inputContent');
        
        if (!titleTextarea || !contentTextarea) {
            console.error('Could not find title or content textarea');
            return;
        }

        const title = titleTextarea.value.trim();
        const content = contentTextarea.value.trim();

        if (title === '') {
            showTitleValidationError(titleTextarea);
            return;
        }

        try {
            const saveButton = noteElement.querySelector('#saveNewNote');
            const originalText = saveButton.textContent;
            saveButton.textContent = 'Saving...';
            saveButton.disabled = true;

            const response = await noteAPI.create(title, content);
            
            showNotification('Note created successfully!', 'success');
            
            titleTextarea.value = '';
            contentTextarea.value = '';
            autoResizeTextarea(titleTextarea);
            autoResizeTextarea(contentTextarea);
            
            if (noteElement.classList.contains("expanded") || noteElement.classList.contains("expanding")) {
              closeExpandedNote(noteElement);
            }
            
            setTimeout(() => {
              window.location.reload();
            }, 1000);
            
        } catch (error) {
            console.error('Error saving note:', error);
            showNotification(`Error saving note: ${error.message}`, 'error');
        } finally {
            const saveButton = noteElement.querySelector('#saveNewNote');
            if (saveButton) {
                saveButton.textContent = 'Save';
                saveButton.disabled = false;
            }
        }
    }

    async function handleSaveExistingNote(noteElement) {
        const noteId = noteElement.getAttribute('data-note-id');
        const titleTextarea = noteElement.querySelector('textarea[name="noteTitle"]');
        const contentTextarea = noteElement.querySelector('textarea[name="noteContent"]');
        
        if (!noteId || !titleTextarea || !contentTextarea) {
            console.error('Could not find note ID or textareas');
            return;
        }

        const content = contentTextarea.value.trim();

        try {
            const saveButton = noteElement.querySelector('.saveNote');
            const originalText = saveButton.textContent;
            saveButton.textContent = 'Saving...';
            saveButton.disabled = true;

            await noteAPI.update(noteId, content, 'overwrite');
            
            showNotification('Note updated successfully!', 'success');
            
            if (noteElement.classList.contains("expanded") || noteElement.classList.contains("expanding")) {
                closeExpandedNote(noteElement);
            }
            
        } catch (error) {
            console.error('Error updating note:', error);
            showNotification(`Error updating note: ${error.message}`, 'error');
        } finally {
            const saveButton = noteElement.querySelector('.saveNote');
            if (saveButton) {
                saveButton.textContent = 'Save';
                saveButton.disabled = false;
            }
        }
    }

    async function handleDeleteNote(noteElement) {
        const noteId = noteElement.getAttribute('data-note-id');
        
        if (!noteId) {
            console.error('Could not find note ID');
            return;
        }

        const confirmDelete = confirm('Are you sure you want to delete this note? This action cannot be undone.');
        if (!confirmDelete) {
            return;
        }

        try {
            await noteAPI.delete(noteId);
            showNotification('Note deleted successfully!', 'success');
            deleteNoteCard(noteElement);
        } catch (error) {
            console.error('Error deleting note:', error);
            showNotification(`Error deleting note: ${error.message}`, 'error');
        }
    }

    // Delete note card function
    function deleteNoteCard(noteElement) {
        if (noteElement.classList.contains("expanded") || noteElement.classList.contains("expanding")) {
          closeExpandedNote(noteElement);
        }
        
        noteElement.classList.add('deleting');
        
        setTimeout(() => {
          noteElement.remove();
        }, 400);
    }
    
    // Clear create note form function
    function clearCreateNoteForm(createNoteElement) {
        const titleTextarea = createNoteElement.querySelector('#inputTitle');
        const contentTextarea = createNoteElement.querySelector('#inputContent');
        
        if (titleTextarea) {
          titleTextarea.value = '';
          removeTitleValidationError(titleTextarea);
          autoResizeTextarea(titleTextarea);
        }
        
        if (contentTextarea) {
          contentTextarea.value = '';
          autoResizeTextarea(contentTextarea);
        }
        
        if (createNoteElement.classList.contains("expanded") || createNoteElement.classList.contains("expanding")) {
          closeExpandedNote(createNoteElement);
        }
        
        if (titleTextarea) {
          titleTextarea.focus();
        }
    }

    // Apply auto-resize to all textareas
    function initAutoResize() {
        document.querySelectorAll('textarea').forEach(textarea => {
          if (textarea.value.trim() !== '' || textarea.textContent.trim() !== '') {
            setTimeout(() => {
              autoResizeTextarea(textarea);
            }, 10);
          } else {
            autoResizeTextarea(textarea);
          }
          
          textarea.addEventListener('input', () => autoResizeTextarea(textarea));
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
        const ellipsisBtn = note.querySelector(".fa-ellipsis")?.parentElement;
        const saveBtn = note.querySelector(".saveNote");
        const saveNewBtn = note.querySelector("#saveNewNote, [id^='save']");

        if (closeBtn && !closeBtn.classList.contains("hidden")) {
          closeBtn.classList.add("hidden");
        }

        const textareas = note.querySelectorAll('textarea');
        textareas.forEach(textarea => {
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
              removeTitleValidationError(titleTextarea);
            }
          });
        }

        // Expand note on click
        note.addEventListener("click", (e) => {
          if (
            e.target.closest(".closeNoteBtn") ||
            e.target.closest(".trashBtn") ||
            e.target.closest("button") ||
            e.target.tagName.toLowerCase() === "button" ||
            e.target.closest("saveNote")
          ) return;

          notes.forEach(n => {
            if (n !== note && n.classList.contains("expanded")) {
              closeExpandedNote(n);
            }
          });

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
            
            if (note.querySelector('#inputTitle')) {
              clearCreateNoteForm(note);
              return;
            }
            
            handleDeleteNote(note);
          });
        }

        // Save existing note button click handler
        if (saveBtn) {
          saveBtn.addEventListener("click", (e) => {
            e.stopPropagation();
            handleSaveExistingNote(note);
          });
        }

        // Save new note button click handler
        if (saveNewBtn) {
          saveNewBtn.addEventListener("click", (e) => {
            e.stopPropagation();
            handleSaveNewNote(note);
          });
        }
    });

    // Search functionality
    const searchBar = document.querySelector("#searchBar");
    const searchBtn = document.querySelector(".searchBtn");
    
    if (searchBar && searchBtn) {
        const performSearch = () => {
            const searchTerm = searchBar.value.trim();
            if (!searchTerm) {
                // Show all notes
                const noteCards = document.querySelectorAll('.note-card[data-note-id]');
                noteCards.forEach(card => {
                    card.style.display = '';
                });
                return;
            }

            // Filter visible notes
            const noteCards = document.querySelectorAll('.note-card[data-note-id]');
            let foundCount = 0;
            
            noteCards.forEach(card => {
                const title = card.querySelector('textarea[name="noteTitle"]')?.value?.toLowerCase() || '';
                const content = card.querySelector('textarea[name="noteContent"]')?.value?.toLowerCase() || '';
                const searchLower = searchTerm.toLowerCase();
                
                if (title.includes(searchLower) || content.includes(searchLower)) {
                    card.style.display = '';
                    foundCount++;
                } else {
                    card.style.display = 'none';
                }
            });
            
            showNotification(`Found ${foundCount} notes matching "${searchTerm}"`, 'info');
        };

        searchBtn.addEventListener("click", performSearch);
        searchBar.addEventListener("keypress", (e) => {
            if (e.key === "Enter") {
                performSearch();
            }
        });
    }

    // Function to restructure note for expanded view
    function restructureExpandedNote(note) {
        const titleArea = note.querySelector('textarea[rows="1"]')?.parentElement;
        const contentArea = note.querySelector('textarea[rows="3"]')?.parentElement || note.querySelector('textarea[rows="4"]')?.parentElement;
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
            const minRows = Math.floor((window.innerHeight - 300) / 24);
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
            setTimeout(() => autoResizeTextarea(titleTextarea), 10);
          }
        }

        if (contentArea) {
          contentArea.classList.remove('content-area');
          const contentTextarea = contentArea.querySelector('textarea');
          if (contentTextarea) {
            contentTextarea.classList.remove('content-textarea', 'auto-resize');
            contentTextarea.setAttribute('rows', '4');
            contentTextarea.style.height = 'auto';
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

        note.classList.add('expanding');
        Object.assign(note.style, originalStyles);

        if (closeBtn) closeBtn.classList.remove("hidden");
        if (ellipsisBtn) ellipsisBtn.classList.add("hidden");

        note.offsetHeight;

        requestAnimationFrame(() => {
          note.style.top = '0px';
          note.style.left = '0px';
          note.style.width = '100vw';
          note.style.height = '100vh';
          note.style.margin = '0';
          note.style.padding = '2rem';
          note.style.borderRadius = '0';
        });

        setTimeout(() => {
          note.classList.remove('expanding');
          note.classList.add('expanded');
          
          note.style.cssText = '';
          
          restructureExpandedNote(note);
          
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

        if (!note.classList.contains("expanded") && !note.classList.contains("expanding")) {
          return;
        }

        resetNoteStructure(note);

        note.classList.remove('expanded', 'expanding');
        note.style.cssText = '';
        
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