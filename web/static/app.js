"use strict";

var API = "/api";

// State.
var currentSort = "date";
var currentOrder = "desc";
var editingId = null;
var itemsCache = []; // кеш для Edit (GET /items/:id не существует в роутере)

// ------- Tabs -------
function switchTab(tab) {
    document.querySelectorAll(".tab").forEach(function (t) {
        t.classList.toggle("active", t.dataset.tab === tab);
    });
    document.getElementById("tab-items").classList.toggle("hidden", tab !== "items");
    document.getElementById("tab-analytics").classList.toggle("hidden", tab !== "analytics");

    if (tab === "analytics") {
        loadAnalytics();
    }
}

// ------- Toast -------
function showToast(message, type) {
    var toast = document.createElement("div");
    toast.className = "toast toast-" + type;
    toast.textContent = message;
    document.body.appendChild(toast);
    setTimeout(function () {
        toast.remove();
    }, 3000);
}

// ------- Items CRUD -------
function loadItems() {
    var params = new URLSearchParams();
    var from = document.getElementById("filter-from").value;
    var to = document.getElementById("filter-to").value;
    var category = document.getElementById("filter-category").value;
    var type = document.getElementById("filter-type").value;

    if (from) params.set("from", from);
    if (to) params.set("to", to);
    if (category) params.set("category", category);
    if (type) params.set("type", type);
    params.set("sort_by", currentSort);
    params.set("order", currentOrder);

    fetch(API + "/items?" + params.toString())
        .then(function (res) {
            if (!res.ok) return res.json().then(function (e) { throw new Error(e.error); });
            return res.json();
        })
        .then(function (data) {
            var items = Array.isArray(data) ? data : (data.items || []);
            itemsCache = items;
            renderItems(items);
        })
        .catch(function (err) {
            showToast(err.message, "error");
        });
}

function renderItems(items) {
    var tbody = document.getElementById("items-table-body");
    var empty = document.getElementById("empty-items");

    if (!items || !Array.isArray(items) || items.length === 0) {
        tbody.innerHTML = "";
        empty.classList.remove("hidden");
        return;
    }

    empty.classList.add("hidden");
    tbody.innerHTML = items.map(function (item) {
        var badgeClass = item.type === "income" ? "badge-income" : "badge-expense";
        var dateStr = item.date ? item.date.substring(0, 10) : "";
        return '<tr>' +
            '<td><span class="badge ' + badgeClass + '">' + escapeHtml(item.type) + '</span></td>' +
            '<td>' + Number(item.amount).toFixed(2) + '</td>' +
            '<td>' + escapeHtml(item.category) + '</td>' +
            '<td>' + escapeHtml(item.description || "") + '</td>' +
            '<td>' + dateStr + '</td>' +
            '<td>' +
            '<button class="btn btn-outline btn-sm" onclick="openEditModal(\'' + item.id + '\')">Edit</button> ' +
            '<button class="btn btn-danger btn-sm" onclick="deleteItem(\'' + item.id + '\')">Delete</button>' +
            '</td>' +
            '</tr>';
    }).join("");
}

function saveItem() {
    var body = {
        type: document.getElementById("item-type").value,
        amount: parseFloat(document.getElementById("item-amount").value),
        category: document.getElementById("item-category").value.trim(),
        description: document.getElementById("item-description").value.trim(),
        date: document.getElementById("item-date").value
    };

    if (!body.category || !body.date || isNaN(body.amount) || body.amount <= 0) {
        showToast("Please fill in all required fields (type, amount > 0, category, date).", "error");
        return;
    }

    fetch(API + "/items", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    })
        .then(function (res) {
            if (!res.ok) return res.json().then(function (e) { throw new Error(e.error); });
            return res.json();
        })
        .then(function () {
            showToast("Item created", "success");
            clearForm();
            loadItems();
        })
        .catch(function (err) {
            showToast(err.message, "error");
        });
}

function deleteItem(id) {
    if (!confirm("Delete this item?")) return;

    fetch(API + "/items/" + id, { method: "DELETE" })
        .then(function (res) {
            if (!res.ok) return res.json().then(function (e) { throw new Error(e.error); });
            showToast("Item deleted", "success");
            loadItems();
        })
        .catch(function (err) {
            showToast(err.message, "error");
        });
}

function openEditModal(id) {
    fetch(API + "/items/" + id)
        .then(function (res) {
            if (!res.ok) return res.json().then(function (e) { throw new Error(e.error); });
            return res.json();
        })
        .then(function (item) {
            editingId = id;
            document.getElementById("modal-type").value = item.type;
            document.getElementById("modal-amount").value = item.amount;
            document.getElementById("modal-category").value = item.category;
            document.getElementById("modal-description").value = item.description || "";
            document.getElementById("modal-date").value = item.date ? item.date.substring(0, 10) : "";
            document.getElementById("edit-modal").classList.add("show");
        })
        .catch(function (err) {
            showToast(err.message, "error");
        });
}

function closeModal() {
    editingId = null;
    document.getElementById("edit-modal").classList.remove("show");
}

function updateItem() {
    if (!editingId) return;

    var body = {
        type: document.getElementById("modal-type").value,
        amount: parseFloat(document.getElementById("modal-amount").value),
        category: document.getElementById("modal-category").value.trim(),
        description: document.getElementById("modal-description").value.trim(),
        date: document.getElementById("modal-date").value
    };

    fetch(API + "/items/" + editingId, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    })
        .then(function (res) {
            if (!res.ok) return res.json().then(function (e) { throw new Error(e.error); });
            return res.json();
        })
        .then(function () {
            showToast("Item updated", "success");
            closeModal();
            loadItems();
        })
        .catch(function (err) {
            showToast(err.message, "error");
        });
}

function clearForm() {
    document.getElementById("item-type").value = "income";
    document.getElementById("item-amount").value = "";
    document.getElementById("item-category").value = "";
    document.getElementById("item-description").value = "";
    document.getElementById("item-date").value = todayStr();
}

function cancelEdit() {
    clearForm();
}

// ------- Sorting -------
function sortBy(field) {
    if (currentSort === field) {
        currentOrder = currentOrder === "asc" ? "desc" : "asc";
    } else {
        currentSort = field;
        currentOrder = "asc";
    }
    loadItems();
}

// ------- Analytics -------
function loadAnalytics() {
    var from = document.getElementById("analytics-from").value;
    var to = document.getElementById("analytics-to").value;
    var groupBy = document.getElementById("analytics-group").value;
    var type = document.getElementById("analytics-type").value;


    if (!from || !to) {
        showToast("Please select 'from' and 'to' dates for analytics.", "error");
        return;
    }

    var params = new URLSearchParams();
    params.set("from", from);
    params.set("to", to);
    if (groupBy) params.set("group_by", groupBy);
    if (type) params.set("type", type);

    fetch(API + "/analytics?" + params.toString())
        .then(function (res) {
            if (!res.ok) return res.json().then(function (e) { throw new Error(e.error); });
            return res.json();
        })
        .then(function (data) {
            renderAnalytics(data);
        })
        .catch(function (err) {
            showToast(err.message, "error");
        });
}

function renderAnalytics(data) {
    document.getElementById("stat-count").textContent = data.count;
    document.getElementById("stat-sum").textContent = Number(data.total_sum).toFixed(2);
    document.getElementById("stat-avg").textContent = Number(data.avg).toFixed(2);
    document.getElementById("stat-median").textContent = Number(data.median).toFixed(2);
    document.getElementById("stat-p90").textContent = Number(data.p90).toFixed(2);

    var groupsCard = document.getElementById("groups-card");
    var groupsTbody = document.getElementById("groups-table-body");

    if (data.groups && data.groups.length > 0) {
        groupsCard.classList.remove("hidden");
        groupsTbody.innerHTML = data.groups.map(function (g) {
            return '<tr>' +
                '<td>' + escapeHtml(g.key) + '</td>' +
                '<td>' + g.count + '</td>' +
                '<td>' + Number(g.total_sum).toFixed(2) + '</td>' +
                '<td>' + Number(g.avg).toFixed(2) + '</td>' +
                '<td>' + Number(g.median).toFixed(2) + '</td>' +
                '<td>' + Number(g.p90).toFixed(2) + '</td>' +
                '</tr>';
        }).join("");
    } else {
        groupsCard.classList.add("hidden");
        groupsTbody.innerHTML = "";
    }
}

// ------- CSV Export -------
function exportCSV() {
    var params = new URLSearchParams();
    var from = document.getElementById("filter-from").value;
    var to = document.getElementById("filter-to").value;
    var category = document.getElementById("filter-category").value;
    var type = document.getElementById("filter-type").value;

    if (from) params.set("from", from);
    if (to) params.set("to", to);
    if (category) params.set("category", category);
    if (type) params.set("type", type);

    window.location.href = API + "/export/csv?" + params.toString();
}

// ------- Helpers -------
function escapeHtml(text) {
    var div = document.createElement("div");
    div.appendChild(document.createTextNode(text));
    return div.innerHTML;
}

function todayStr() {
    var d = new Date();
    return d.getFullYear() + "-" +
        String(d.getMonth() + 1).padStart(2, "0") + "-" +
        String(d.getDate()).padStart(2, "0");
}

// ------- Init -------
(function init() {
    document.getElementById("item-date").value = todayStr();

    var now = new Date();
    var thirtyDaysAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
    document.getElementById("analytics-from").value = thirtyDaysAgo.getFullYear() + "-" +
        String(thirtyDaysAgo.getMonth() + 1).padStart(2, "0") + "-" +
        String(thirtyDaysAgo.getDate()).padStart(2, "0");
    document.getElementById("analytics-to").value = todayStr();

    loadItems();
})();