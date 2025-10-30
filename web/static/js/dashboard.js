(() => {
  const root = document.getElementById('dashboard-root');
  let activeStatus = '';

  const getReloadInterval = () => {
    if (!root) return 0;
    const value = Number(root.dataset.reloadAfter || 0);
    return Number.isNaN(value) ? 0 : value;
  };

  const setAlert = (type, message) => {
    const el = document.getElementById('jobs-console-alert');
    if (!el) return;
    if (!message) {
      el.classList.add('d-none');
      el.textContent = '';
      return;
    }
    el.textContent = message;
    el.className = `alert alert-${type}`;
  };

  const statusColor = (status) => {
    switch ((status || '').toLowerCase()) {
      case 'pending':
        return 'warning';
      case 'retry':
        return 'info';
      case 'running':
        return 'primary';
      case 'failed':
        return 'danger';
      case 'done':
        return 'success';
      default:
        return 'secondary';
    }
  };

  const renderJobs = (jobs) => {
    const table = document.querySelector('#jobs-console-table tbody');
    if (!table) return;
    if (!jobs || jobs.length === 0) {
      table.innerHTML = '<tr><td colspan="5" class="text-center text-muted py-4">No jobs found for the selected filter.</td></tr>';
      return;
    }

    const rows = jobs
      .map((job) => {
        const attempts = `${job.attempts}/${job.max_attempts}`;
        const runAt = new Date(job.run_at).toLocaleString();
        const actions = [];
        if (['failed', 'retry'].includes(job.status)) {
          actions.push(`<button type="button" class="btn btn-sm btn-outline-primary jobs-action" data-action="retry" data-id="${job.id}">Retry</button>`);
        }
        if (['pending', 'retry', 'running'].includes(job.status)) {
          actions.push(`<button type="button" class="btn btn-sm btn-outline-danger jobs-action" data-action="cancel" data-id="${job.id}">Cancel</button>`);
        }
        const actionsHtml = actions.length ? actions.join(' ') : '<span class="text-muted small">—</span>';
        return `
          <tr>
            <td class="fw-medium">${job.name}</td>
            <td><span class="badge text-bg-${statusColor(job.status)} text-capitalize">${job.status}</span></td>
            <td>${attempts}</td>
            <td>${runAt}</td>
            <td>${actionsHtml}</td>
          </tr>`;
      })
      .join('');

    table.innerHTML = rows;
  };

  const fetchJobs = async (status) => {
    try {
      const params = new URLSearchParams();
      if (status) params.set('status', status);
      params.set('limit', '50');
      const response = await fetch(`/jobs?${params.toString()}`);
      if (!response.ok) {
        throw new Error(`Jobs request failed: ${response.status}`);
      }
      const data = await response.json();
      renderJobs(data);
      setAlert(null, null);
    } catch (err) {
      console.warn('jobs fetch failed', err);
      setAlert('danger', err.message || 'Unable to load jobs.');
    }
  };

  const postJobAction = async (id, action) => {
    try {
      const endpoint = action === 'retry' ? `/jobs/${id}/retry` : `/jobs/${id}/cancel`;
      const body = action === 'cancel' ? JSON.stringify({ reason: 'Cancelled via console' }) : undefined;
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: body ? { 'Content-Type': 'application/json' } : undefined,
        body,
      });
      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Action failed: ${response.status}`);
      }
      setAlert('success', action === 'retry' ? 'Job re-queued.' : 'Job cancelled.');
      await fetchJobs(activeStatus);
    } catch (err) {
      console.warn('job action failed', err);
      setAlert('danger', err.message || 'Job action failed.');
    }
  };

  const bindFilters = () => {
    const buttons = document.querySelectorAll('.jobs-filter');
    buttons.forEach((button) => {
      button.addEventListener('click', () => {
        buttons.forEach((btn) => btn.classList.remove('active'));
        button.classList.add('active');
        activeStatus = button.dataset.status || '';
        fetchJobs(activeStatus);
      });
    });

    const table = document.querySelector('#jobs-console-table');
    if (table) {
      table.addEventListener('click', (event) => {
        const target = event.target;
        if (!(target instanceof HTMLElement)) return;
        if (!target.classList.contains('jobs-action')) return;
        const jobId = target.dataset.id;
        const action = target.dataset.action;
        if (!jobId || !action) return;
        postJobAction(jobId, action);
      });
    }
  };

  const bindJobForm = () => {
    const form = document.getElementById('job-create-form');
    if (!form) return;

    const successAlert = document.getElementById('jobSuccessAlert');
    const errorAlert = document.getElementById('jobErrorAlert');

    const showSuccess = (message) => {
      if (successAlert) {
        successAlert.textContent = message;
        successAlert.classList.remove('d-none');
      }
      if (errorAlert) {
        errorAlert.classList.add('d-none');
      }
    };

    const showError = (message) => {
      if (errorAlert) {
        errorAlert.textContent = message;
        errorAlert.classList.remove('d-none');
      }
      if (successAlert) {
        successAlert.classList.add('d-none');
      }
    };

    form.addEventListener('submit', async (event) => {
      event.preventDefault();

      const formData = new FormData(form);
      const name = formData.get('name');
      const runAtRaw = formData.get('run_at');
      const maxAttempts = Number(formData.get('max_attempts')) || 1;
      const payloadRaw = formData.get('payload');

      if (!name) {
        showError('Job name is required.');
        return;
      }

      let payload;
      if (payloadRaw && payloadRaw.trim().length > 0) {
        try {
          payload = JSON.parse(payloadRaw);
        } catch (err) {
          showError('Payload must be valid JSON.');
          return;
        }
      } else {
        payload = {};
      }

      const body = {
        name,
        payload,
        max_attempts: maxAttempts,
      };

      if (runAtRaw) {
        const iso = new Date(runAtRaw).toISOString();
        if (Number.isNaN(Date.parse(iso))) {
          showError('Run At must be a valid date.');
          return;
        }
        body.run_at = iso;
      }

      try {
        const response = await fetch('/jobs', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(body),
        });

        if (!response.ok) {
          const text = await response.text();
          throw new Error(text || response.statusText);
        }

        showSuccess('Job queued successfully. Refreshing…');
        const modalElement = document.getElementById('jobCreateModal');
        if (modalElement && window.bootstrap) {
          const modalInstance = bootstrap.Modal.getInstance(modalElement) || new bootstrap.Modal(modalElement);
          modalInstance.hide();
        }

        setTimeout(() => {
          window.location.reload();
        }, 800);
      } catch (err) {
        showError(err.message || 'Unable to queue job.');
      }
    });
  };

  const initialize = () => {
    bindJobForm();
    bindFilters();
    fetchJobs(activeStatus);
  };

  const setupAutoRefresh = () => {
    const interval = getReloadInterval();
    initialize();

    const refresh = async () => {
      try {
        const response = await fetch('/partials/dashboard', {
          headers: { 'Accept': 'text/html' },
        });
        if (!response.ok) {
          throw new Error(`refresh failed: ${response.status}`);
        }
        const html = await response.text();
        const parser = new DOMParser();
        const doc = parser.parseFromString(html, 'text/html');
        const replacement = doc.querySelector('#dashboard-root');
        if (replacement && root) {
          root.innerHTML = replacement.innerHTML;
          const reloadAttr = replacement.dataset.reloadAfter;
          if (reloadAttr) {
            root.dataset.reloadAfter = reloadAttr;
          }
          initialize();
        }
      } catch (err) {
        console.warn('auto refresh failed', err);
      }
    };

    window.addEventListener('shanraq:refresh', refresh);
    if (interval) {
      setInterval(refresh, interval * 1000);
    }
  };

  setupAutoRefresh();
})();
