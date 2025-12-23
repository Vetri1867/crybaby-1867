
import { ref, uploadBytes, getDownloadURL, listAll, deleteObject } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-storage.js";
import { collection, addDoc, getDocs, deleteDoc, doc, query, where } from "https://www.gstatic.com/firebasejs/10.12.2/firebase-firestore.js";
import { storage, db, auth } from "./firebase.js";

const bookFileInput = document.getElementById("bookFile");
const bookSubjectInput = document.getElementById("bookSubject");

export async function uploadBook() {
    const file = bookFileInput.files[0];
    const subject = bookSubjectInput.value.trim();
    const user = auth.currentUser;

    if (!file || !subject || !user) {
        alert("Please select a file and enter a subject.");
        return;
    }

    const storageRef = ref(storage, `books/${user.uid}/${subject}/${file.name}`);

    try {
        await uploadBytes(storageRef, file);
        const downloadURL = await getDownloadURL(storageRef);

        await addDoc(collection(db, "books"), {
            name: file.name,
            subject,
            url: downloadURL,
            userId: user.uid,
            createdAt: new Date()
        });

        alert("Book uploaded successfully!");
        bookFileInput.value = "";
        bookSubjectInput.value = "";
        loadBooks();
    } catch (error) {
        console.error("Error uploading book: ", error);
        alert("Failed to upload book.");
    }
}

export async function loadBooks() {
    const user = auth.currentUser;
    if (!user) return;

    const booksListParent = document.getElementById("booksListParent");
    const booksListStudent = document.getElementById("booksListStudent");

    booksListParent.innerHTML = "";
    booksListStudent.innerHTML = "";

    const q = query(collection(db, "books"), where("userId", "==", user.uid));
    const querySnapshot = await getDocs(q);

    querySnapshot.forEach((doc) => {
        const book = doc.data();
        const bookId = doc.id;

        const parentItem = document.createElement("div");
        parentItem.innerHTML = `<p>${book.name} (${book.subject}) <button data-id="${bookId}" data-url="${book.url}">Delete</button></p>`;
        booksListParent.appendChild(parentItem);

        const studentItem = document.createElement("div");
        studentItem.innerHTML = `<p><a href="${book.url}" target="_blank">${book.name}</a> (${book.subject})</p>`;
        booksListStudent.appendChild(studentItem);
    });

    document.querySelectorAll("#booksListParent button").forEach(button => {
        button.addEventListener("click", async (e) => {
            const bookId = e.target.dataset.id;
            const bookUrl = e.target.dataset.url;
            deleteBook(bookId, bookUrl);
        });
    });
}

async function deleteBook(bookId, bookUrl) {
    const user = auth.currentUser;
    if (!user) return;

    const bookRef = doc(db, "books", bookId);
    const storageRef = ref(storage, bookUrl);

    try {
        await deleteDoc(bookRef);
        await deleteObject(storageRef);
        alert("Book deleted successfully!");
        loadBooks();
    } catch (error) {
        console.error("Error deleting book: ", error);
        alert("Failed to delete book.");
    }
}
