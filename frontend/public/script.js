// Load API integration first
window.addEventListener("load", function () {
  // Dev
  //const API_BASE_URL = '/api/v1';
  // Prod
  const API_BASE_URL = 'https://andika.adamndegwa.workers.dev'

  const apiCall = async (endpoint, options = {}) => {
    const url = `${API_BASE_URL}${endpoint}`;
    const config = {
      headers: { 'Content-Type': 'application/json' },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      const contentType = response.headers.get('content-type') || '';
      if (contentType.includes('application/json')) {
        const data = await response.json();
        if (!response.ok) throw new Error(data.error || `HTTP error! status: ${response.status}`);
        return data;
      } else {
        const text = await response.text();
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status} - ${text}`);
        return text;
      }
    } catch (error) {
      console.error('API call failed:', error);
      throw error;
    }
  };

  const noteAPI = {
    create: async (name, content) => apiCall('/notes', { method: 'POST', body: JSON.stringify({ name, content }) }),
    list: async () => apiCall('/notes'),
    get: async (id) => apiCall(`/notes/${id}`),
    update: async (id, content, mode = 'overwrite') => apiCall(`/notes/${id}`, { method: 'PUT', body: JSON.stringify({ mode, content }) }),
    delete: async (id) => apiCall(`/notes/${id}`, { method: 'DELETE' }),
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
    setTimeout(() => { notification.style.opacity = '0'; setTimeout(() => notification.remove(), 300); }, 3000);
  };

  // Header Menu toggles
  const showMenuBtn = document.querySelector("#showMenu");
  const hideMenuBtn = document.querySelector("#hideMenu");
  const mobileNav = document.querySelector("#mobileNav");
  if (showMenuBtn && hideMenuBtn && mobileNav) {
    showMenuBtn.addEventListener("click", () => { mobileNav.classList.remove("hidden", "animate-fade-out"); mobileNav.classList.add("animate-fade-in"); });
    hideMenuBtn.addEventListener("click", () => {
      mobileNav.classList.remove("animate-fade-in");
      mobileNav.classList.add("animate-fade-out");
      mobileNav.addEventListener("animationend", () => { if (mobileNav.classList.contains("animate-fade-out")) mobileNav.classList.add("hidden"); }, { once: true });
    });
  }

  // FAQ section
  document.querySelectorAll("[toggleElement]").forEach(toggle => {
    toggle.addEventListener("click", function () {
      const answerElement = toggle.querySelector("[answer]");
      const caretElement = toggle.querySelector("img");
      if (answerElement.classList.contains("hidden")) {
        answerElement.classList.remove("hidden"); caretElement.classList.add("rotate-90");
      } else { answerElement.classList.add("hidden"); caretElement.classList.remove("rotate-90"); }
    });
  });

  // Loading screen
  const loader = document.getElementById("loadingScreen");
  const app = document.getElementById("app");
  const appNotes = document.getElementById("appNotes");
  if (loader) {
    setTimeout(() => {
      loader.classList.add("animate-fade-out");
      loader.addEventListener("animationend", () => {
        loader.style.display = "none";
        if (app) { app.classList.remove("hidden"); app.classList.add("animate-fade-in"); }
        else if (appNotes) { appNotes.classList.remove("hidden"); appNotes.classList.add("animate-fade-in"); setTimeout(() => { initAutoResize(); }, 100); }
      }, { once: true });
    }, 2800);
  } else { setTimeout(() => { initAutoResize(); }, 100); }

  // Auto-resize textarea function
  function autoResizeTextarea(textarea) {
    const originalHeight = textarea.style.height;
    textarea.style.height = 'auto';
    let newHeight = textarea.scrollHeight;
    if (textarea.getAttribute('rows') === '1') newHeight = Math.max(newHeight, parseInt(getComputedStyle(textarea).lineHeight) || 24);
    textarea.style.height = newHeight + 'px';
    if (Math.abs(newHeight - parseInt(originalHeight || '0')) > 5) textarea.offsetHeight;
  }

  // Validation helper functions
  function showTitleValidationError(titleTextarea) {
    const existingError = titleTextarea.parentElement.querySelector('.title-error-message');
    if (existingError) return;
    const errorMessage = document.createElement('div');
    errorMessage.className = 'title-error-message';
    errorMessage.textContent = 'A title is necessary in order to save';
    titleTextarea.parentElement.insertBefore(errorMessage, titleTextarea.nextSibling);
    titleTextarea.classList.add('error-state');
    titleTextarea.focus();
  }
  function removeTitleValidationError(titleTextarea) {
    const errorMessage = titleTextarea.parentElement.querySelector('.title-error-message');
    if (errorMessage) errorMessage.remove();
    titleTextarea.classList.remove('error-state');
  }

  // API-integrated functions (used by Save/Delete buttons)
  async function handleSaveExistingNote(noteElement) {
    const noteId = noteElement.getAttribute('data-note-id');
    const contentTextarea = noteElement.querySelector('textarea[name="noteContent"]');
    if (!noteId || !contentTextarea) return;
    const content = contentTextarea.value.trim();
    try {
      const saveButton = noteElement.querySelector('.saveNote');
      saveButton.textContent = 'Saving...';
      saveButton.disabled = true;
      await noteAPI.update(noteId, content, 'overwrite');
      showNotification('Note updated successfully!', 'success');
      resetAllNotesToDefault();
    } catch (error) { console.error('Error updating note:', error); showNotification(`Error updating note: ${error.message}`, 'error'); }
    finally { const saveButton = noteElement.querySelector('.saveNote'); if (saveButton) { saveButton.textContent = 'Save'; saveButton.disabled = false; } }
  }

  async function handleDeleteNote(noteElement) {
    const noteId = noteElement.getAttribute('data-note-id');
    if (!noteId) return;
    const confirmDelete = confirm('Are you sure you want to delete this note? This action cannot be undone.');
    if (!confirmDelete) return;
    try { await noteAPI.delete(noteId); showNotification('Note deleted successfully!', 'success'); deleteNoteCard(noteElement); } 
    catch (error) { console.error('Error deleting note:', error); showNotification(`Error deleting note: ${error.message}`, 'error'); }
  }
  function deleteNoteCard(noteElement) { noteElement.classList.add('deleting'); setTimeout(() => { noteElement.remove(); }, 400); }

  // Apply auto-resize to all textareas
  function initAutoResize() {
    document.querySelectorAll('textarea').forEach(textarea => {
      if (textarea.value.trim() !== '' || textarea.textContent.trim() !== '') setTimeout(() => { autoResizeTextarea(textarea); }, 10);
      else autoResizeTextarea(textarea);
      textarea.addEventListener('input', () => autoResizeTextarea(textarea));
      textarea.addEventListener('focus', () => { setTimeout(() => autoResizeTextarea(textarea), 10); });
    });
  }

  // Note expansion
  function attachNoteHandlers(note) {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const trashBtn = note.querySelector(".trashBtn");
    const saveBtn = note.querySelector(".saveNote");
    if (closeBtn && !closeBtn.classList.contains("hidden")) closeBtn.classList.add("hidden");
    note.querySelectorAll('textarea').forEach(textarea => setTimeout(() => autoResizeTextarea(textarea), 50));

    note.addEventListener("click", e => {
      if (e.target.closest(".closeNoteBtn") || e.target.closest(".trashBtn") || e.target.closest("button") || e.target.tagName.toLowerCase() === "button") return;
      document.querySelectorAll(".note-card").forEach(n => { if (n !== note && n.classList.contains("expanded")) closeExpandedNote(n); });
      expandNote(note);
    });

    if (closeBtn) closeBtn.addEventListener("click", e => { e.stopPropagation(); closeExpandedNote(note); });
    if (trashBtn) trashBtn.addEventListener("click", e => { e.stopPropagation(); handleDeleteNote(note); });
    if (saveBtn) saveBtn.addEventListener("click", e => { e.stopPropagation(); handleSaveExistingNote(note); });
  }

  document.querySelectorAll(".note-card").forEach(note => attachNoteHandlers(note));

  document.body.addEventListener('htmx:afterSwap', (evt) => {
    const swapped = evt.detail?.elt;
    if (swapped) {
      swapped.querySelectorAll?.('.note-card')?.forEach(n => attachNoteHandlers(n));
      swapped.querySelectorAll?.('textarea')?.forEach(t => { autoResizeTextarea(t); t.addEventListener('input', () => autoResizeTextarea(t)); });
    }
  });

  // Expand/close helpers
  function restructureExpandedNote(note) {
    const titleArea = note.querySelector('textarea[rows="1"]')?.parentElement;
    const contentArea = note.querySelector('textarea[name="noteContent"]')?.parentElement;
    const buttonArea = note.querySelector('.flex.justify-between');
    if (titleArea) { titleArea.classList.add('title-area'); titleArea.querySelector('textarea')?.classList.add('title-textarea'); }
    if (contentArea) {
      contentArea.classList.add('content-area');
      const contentTextarea = contentArea.querySelector('textarea');
      if (contentTextarea) contentTextarea.setAttribute('rows', Math.max(Math.floor((window.innerHeight - 300) / 24), 10));
    }
    if (buttonArea) buttonArea.classList.add('button-area');
  }

  function resetNoteStructure(note) {
    const titleArea = note.querySelector('.title-area');
    const contentArea = note.querySelector('.content-area');
    const buttonArea = note.querySelector('.button-area');
    if (titleArea) { titleArea.classList.remove('title-area'); titleArea.querySelector('textarea')?.classList.remove('title-textarea'); }
    if (contentArea) {
      contentArea.classList.remove('content-area');
      const contentTextarea = contentArea.querySelector('textarea');
      if (contentTextarea) { contentTextarea.classList.remove('content-textarea', 'auto-resize'); contentTextarea.setAttribute('rows', '4'); autoResizeTextarea(contentTextarea); }
    }
    if (buttonArea) buttonArea.classList.remove('button-area');
  }

  function expandNote(note) {
    const closeBtn = note.querySelector(".closeNoteBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis")?.parentElement;
    const rect = note.getBoundingClientRect();
    const originalStyles = { width: rect.width + 'px', height: rect.height + 'px', top: rect.top + 'px', left: rect.left + 'px', padding: getComputedStyle(note).padding, margin: getComputedStyle(note).margin, borderRadius: getComputedStyle(note).borderRadius };
    note.classList.add('expanding'); Object.assign(note.style, originalStyles);
    if (closeBtn) closeBtn.classList.remove("hidden");
    if (ellipsisBtn) ellipsisBtn.classList.add("hidden");
    note.offsetHeight;
    requestAnimationFrame(() => { note.style.top = '0px'; note.style.left = '0px'; note.style.width = '100vw'; note.style.height = '100vh'; note.style.margin = '0'; note.style.padding = '2rem'; note.style.borderRadius = '0'; });
    setTimeout(() => { note.classList.remove('expanding'); note.classList.add('expanded'); note.style.cssText = ''; restructureExpandedNote(note); note.querySelectorAll('textarea').forEach(t => setTimeout(() => autoResizeTextarea(t), 50)); }, 500);
  }

  function closeExpandedNote(note) {
    if (!note.classList.contains("expanded") && !note.classList.contains("expanding")) return;
    resetNoteStructure(note);
    note.classList.remove('expanded', 'expanding');
    note.style.cssText = '';
    const closeBtn = note.querySelector(".closeNoteBtn");
    const ellipsisBtn = note.querySelector(".fa-ellipsis")?.parentElement;
    if (closeBtn) closeBtn.classList.add("hidden");
    if (ellipsisBtn) ellipsisBtn.classList.remove("hidden");
  }

  document.addEventListener("keydown", (e) => { if (e.key === "Escape") document.querySelectorAll(".note-card.expanded,.note-card.expanding").forEach(note => closeExpandedNote(note)); });
  window.addEventListener("resize", () => {
    document.querySelectorAll(".note-card").forEach(note => {
      if (note.classList.contains("expanded")) {
        const contentTextarea = note.querySelector('.content-textarea');
        if (contentTextarea) { contentTextarea.setAttribute('rows', Math.max(Math.floor((window.innerHeight - 300) / 24), 10)); autoResizeTextarea(contentTextarea); }
      } else if (note.classList.contains("expanding")) { note.style.width = '100vw'; note.style.height = '100vh'; }
    });
  });

  // === NEW: reset all notes and Create Note form to default size ===
  function resetAllNotesToDefault() {
    document.querySelectorAll('.note-card.expanded, .note-card.expanding').forEach(note => closeExpandedNote(note));
    document.querySelectorAll('.note-card textarea').forEach(textarea => {
      textarea.style.height = 'auto';
      if (textarea.name === 'noteContent') textarea.setAttribute('rows', '4');
      else if (textarea.rows === 1) textarea.setAttribute('rows', '1');
      autoResizeTextarea(textarea);
    });
    const createForm = document.getElementById('create-note-form');
    if (createForm) {
      createForm.reset();
      createForm.querySelectorAll('textarea').forEach(textarea => {
        textarea.style.height = 'auto';
        if (textarea.id === 'inputContent') textarea.setAttribute('rows', '4');
        if (textarea.id === 'inputTitle') textarea.setAttribute('rows', '1');
        autoResizeTextarea(textarea);
      });
    }
  }

  // HTMX afterRequest for Create Note form
  document.body.addEventListener('htmx:afterRequest', (evt) => {
    const target = evt.target || evt.detail?.elt;
    if (target && target.id === 'create-note-form') {
      resetAllNotesToDefault();
    }
  });

});
