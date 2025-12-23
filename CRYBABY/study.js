import { db } from "./firebase.js";
import { addDoc, collection } from
"https://www.gstatic.com/firebasejs/10.12.2/firebase-firestore.js";

export async function logStudy(userId, subject, duration) {
  await addDoc(collection(db, "studyLogs"), {
    userId,
    subject,
    duration,
    date: new Date()
  });
}
