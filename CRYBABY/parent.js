import { db } from "./firebase.js";
import { getDocs, query, where, collection } from
"https://www.gstatic.com/firebasejs/10.12.2/firebase-firestore.js";

export async function getChildProgress(childId) {
  const q = query(
    collection(db, "studyLogs"),
    where("userId", "==", childId)
  );

  const snap = await getDocs(q);
  return snap.docs.map(d => d.data());
}
