
import {
    onAuthStateChanged,
    createUserWithEmailAndPassword,
    signInWithEmailAndPassword,
    signOut
} from "https://www.gstatic.com/firebasejs/10.12.2/firebase-auth.js";
import { doc, setDoc, getDoc } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-firestore.js";
import { auth, db } from "./firebase.js";
import { showStudentDashboard, showParentDashboard, showAuthCard, showLoader, hideLoader } from "./ui.js";

const emailInput = document.getElementById("email");
const passwordInput = document.getElementById("password");
const roleSelect = document.getElementById("role");

export function initAuth(appContent, authCard, studentPanel, parentPanel) {
    onAuthStateChanged(auth, async (user) => {
        hideLoader();
        appContent.classList.remove("hidden");

        authCard.classList.add("hidden");
        studentPanel.classList.add("hidden");
        parentPanel.classList.add("hidden");

        if (!user) {
            showAuthCard();
            return;
        }

        try {
            const snap = await getDoc(doc(db, "users", user.uid));
            const role = snap.exists() ? snap.data().role : null;

            if (role === "parent") {
                showParentDashboard();
            } else {
                showStudentDashboard();
            }
        } catch {
            showAuthCard();
        }
    });
}

export async function signupUser() {
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    const role = roleSelect.value;

    if (!email || !password || password.length < 6) {
        alert("Enter valid email & password (6+ chars)");
        return;
    }

    try {
        showLoader();
        const userCred = await createUserWithEmailAndPassword(auth, email, password);
        await setDoc(doc(db, "users", userCred.user.uid), {
            email,
            role,
            createdAt: new Date()
        });
        alert("Signup successful");
    } catch (e) {
        alert(e.message);
    } finally {
        hideLoader();
    }
}

export async function loginUser() {
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();

    if (!email || !password) {
        alert("Please enter email and password");
        return;
    }

    try {
        showLoader();
        await signInWithEmailAndPassword(auth, email, password);
    } catch (e) {
        alert(e.message);
    } finally {
        hideLoader();
    }
}

export async function logoutUser() {
    try {
        await signOut(auth);
    } catch (e) {
        alert(e.message);
    }
}
