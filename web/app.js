const API_KEY = 'default-api-key';

async function apiGet(path) {
  const res = await fetch(path, { headers: { 'X-Api-Key': API_KEY } });
  return res.json();
}

async function apiPost(path, body) {
  const res = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'X-Api-Key': API_KEY },
    body: JSON.stringify(body)
  });
  return res.json();
}

async function loadStats() {
  try {
    const data = await apiGet('/api/stats/overview');
    if (data.code === 0 && data.data) {
      document.getElementById('biz-count').textContent = data.data.biz_count || 0;
      document.getElementById('active-nodes').textContent = data.data.active_nodes || 0;
      document.getElementById('today-allocated').textContent = data.data.today_allocated || 0;
      document.getElementById('peak-qps').textContent = (data.data.peak_qps || 0).toFixed(1);
    }
  } catch (e) {
    console.error('加载统计失败', e);
  }
}

async function loadBizTypes() {
  try {
    const data = await apiGet('/api/biz-types');
    const tbody = document.querySelector('#biz-table tbody');
    tbody.innerHTML = '';
    const select = document.getElementById('biz-select');
    select.innerHTML = '<option value="">选择业务类型</option>';
    if (data.code === 0 && data.data && data.data.items) {
      data.data.items.forEach(b => {
        const tr = document.createElement('tr');
        const statusClass = b.enabled ? 'status-enabled' : 'status-disabled';
        const statusText = b.enabled ? '启用' : '禁用';
        tr.innerHTML = `
          <td>${b.id.substring(0, 8)}...</td>
          <td>${b.name}</td>
          <td>${b.code}</td>
          <td>${b.mode}</td>
          <td class="${statusClass}">${statusText}</td>
          <td>${new Date(b.created_at).toLocaleString()}</td>
        `;
        tbody.appendChild(tr);

        const opt = document.createElement('option');
        opt.value = b.id;
        opt.textContent = `${b.name} (${b.mode})`;
        select.appendChild(opt);
      });
    }
  } catch (e) {
    console.error('加载业务类型失败', e);
  }
}

function showResult(text) {
  document.getElementById('generate-result').textContent = text;
}

async function generateID() {
  const bizID = document.getElementById('biz-select').value;
  if (!bizID) {
    showResult('请先选择业务类型');
    return;
  }
  try {
    const data = await apiPost('/api/id-gen/batch', { biz_type_id: bizID, batch_size: 1 });
    if (data.code === 0 && data.data) {
      showResult(`生成成功！\nID: ${data.data.ids[0]}\n模式: ${data.data.mode}\n业务: ${data.data.biz_type}`);
    } else {
      showResult('生成失败: ' + (data.message || '未知错误'));
    }
  } catch (e) {
    showResult('请求失败: ' + e.message);
  }
}

async function generateBatch() {
  const bizID = document.getElementById('biz-select').value;
  if (!bizID) {
    showResult('请先选择业务类型');
    return;
  }
  try {
    const data = await apiPost('/api/id-gen/batch', { biz_type_id: bizID, batch_size: 10 });
    if (data.code === 0 && data.data) {
      showResult(`批量生成成功！\n共 ${data.data.count} 个 ID\n模式: ${data.data.mode}\n业务: ${data.data.biz_type}\n\nIDs:\n${data.data.ids.join('\n')}`);
    } else {
      showResult('生成失败: ' + (data.message || '未知错误'));
    }
  } catch (e) {
    showResult('请求失败: ' + e.message);
  }
}

document.getElementById('btn-generate').addEventListener('click', generateID);
document.getElementById('btn-batch').addEventListener('click', generateBatch);

loadStats();
loadBizTypes();
