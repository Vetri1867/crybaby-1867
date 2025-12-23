import { db } from "./firebase.js";
import { doc, updateDoc, increment } from
"https://www.gstatic.com/firebasejs/10.12.2/firebase-firestore.js";

export async function addStar(userId) {
  await updateDoc(doc(db, "users", userId), {
    stars: increment(1)
  });
}
