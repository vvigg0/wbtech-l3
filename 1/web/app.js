const notifyForm = document.getElementById("notify-form");
const statusForm = document.getElementById("status-form");
const deleteForm = document.getElementById("delete-form");

const addNotificationBtn = document.getElementById("add-notification-btn");
const notificationsList = document.getElementById("notifications-list");
const notifTemplate = document.getElementById("notif-template");

const statusResult = document.getElementById("status-result");
const deleteResult = document.getElementById("delete-result");
const serverResponse = document.getElementById("server-response");

function showResponse(data) {
  serverResponse.textContent =
    typeof data === "string" ? data : JSON.stringify(data, null, 2);
}

function showBlock(el, data) {
  el.classList.remove("hidden");
  el.textContent =
    typeof data === "string" ? data : JSON.stringify(data, null, 2);
}

function toRFC3339(datetimeLocalValue) {
  if (!datetimeLocalValue) return "";
  return new Date(datetimeLocalValue).toISOString();
}

addNotificationBtn.addEventListener("click", () => {
  const clone = notifTemplate.content.cloneNode(true);
  notificationsList.appendChild(clone);
});

notificationsList.addEventListener("click", (e) => {
  if (e.target.classList.contains("remove-btn")) {
    const items = notificationsList.querySelectorAll(".notif-item");
    if (items.length === 1) {
      alert("Хотя бы один блок уведомления должен остаться.");
      return;
    }
    e.target.closest(".notif-item").remove();
  }
});

notifyForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const items = notificationsList.querySelectorAll(".notif-item");
  const notifs = [];

  for (const item of items) {
    const text = item.querySelector('[name="text"]').value.trim();
    const telegramIDRaw = item.querySelector('[name="telegram_ID"]').value.trim();
    const sendAt = item.querySelector('[name="send_at"]').value;

    if (!/^\d+$/.test(telegramIDRaw)) {
      alert("Telegram ID должен содержать только цифры");
      return;
    }

    if (!sendAt) {
      alert("Укажи дату и время отправки");
      return;
    }

    notifs.push({
      text: text,
      telegram_ID: Number(telegramIDRaw),
      send_at: toRFC3339(sendAt),
    });
  }

  const payload = { notifs };

  try {
    const res = await fetch("/notify", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await res.json().catch(() => ({ error: "invalid json response" }));
    showResponse(data);

    if (!res.ok) {
      alert("Ошибка при создании уведомления");
      return;
    }

    alert("Уведомление(я) отправлено(ы) на сервер");
  } catch (err) {
    showResponse(String(err));
    alert("Ошибка сети");
  }
});

statusForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const id = document.getElementById("status-id").value.trim();
  if (!id) return;

  try {
    const res = await fetch(`/notify/${encodeURIComponent(id)}`);
    const data = await res.json().catch(() => ({ error: "invalid json response" }));

    showResponse(data);
    showBlock(statusResult, data);
  } catch (err) {
    showResponse(String(err));
    showBlock(statusResult, String(err));
  }
});

deleteForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const id = document.getElementById("delete-id").value.trim();
  if (!id) return;

  try {
    const res = await fetch(`/notify/${encodeURIComponent(id)}`, {
      method: "DELETE",
    });

    const data = await res.json().catch(() => ({ error: "invalid json response" }));

    showResponse(data);
    showBlock(deleteResult, data);
  } catch (err) {
    showResponse(String(err));
    showBlock(deleteResult, String(err));
  }
});