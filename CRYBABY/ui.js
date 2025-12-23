
const loader = document.getElementById("loader");
const appContent = document.getElementById("appContent");
const authCard = document.getElementById("authCard");
const studentPanel = document.getElementById("studentPanel");
const parentPanel = document.getElementById("parentPanel");

export function showLoader() {
    loader.style.display = "flex";
}

export function hideLoader() {
    loader.style.display = "none";
}

export function showApp() {
    appContent.classList.remove("hidden");
}

export function showAuthCard() {
    authCard.classList.remove("hidden");
    studentPanel.classList.add("hidden");
    parentPanel.classList.add("hidden");
}

export function showStudentDashboard() {
    studentPanel.classList.remove("hidden");
    authCard.classList.add("hidden");
    parentPanel.classList.add("hidden");
}

export function showParentDashboard() {
    parentPanel.classList.remove("hidden");
    authCard.classList.add("hidden");
    studentPanel.classList.add("hidden");
}
