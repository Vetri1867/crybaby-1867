
import { initializeApp } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-app.js";
import { getAuth } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-auth.js";
import { getFirestore } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-firestore.js";
import { getStorage } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-storage.js";

const firebaseConfig = {
    apiKey: "AIzaSyAOTqlPLO7avWXrZ0JKMSxHSkT5UpS-CGk",
    authDomain: "crybaby-1867.firebaseapp.com",
    projectId: "crybaby-1867",
    storageBucket: "crybaby-1867.appspot.com",
    messagingSenderId: "190342973304",
    appId: "1:190342973304:web:acdef1d189ba367a2b2a8d"
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const db = getFirestore(app);
const storage = getStorage(app);

export { app, auth, db, storage };
