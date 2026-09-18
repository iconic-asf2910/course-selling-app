const API = "http://localhost:8080";

let userToken = "";

async function login() {

    const email = document.getElementById("loginEmail").value;
    const password = document.getElementById("loginPassword").value;

    const response = await fetch(`${API}/user/login`, {
        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify({
            email: email,
            password: password
        })
    });

    const data = await response.json();

    if (response.ok) {
        userToken = data.token;

        document.getElementById("message").innerText =
            "Login successful!";
    } else {
        document.getElementById("message").innerText =
            data.message || "Login failed";
    }
}


async function getCourses() {

    const response = await fetch(`${API}/courses`);

    const courses = await response.json();

    const container = document.getElementById("courses");

    container.innerHTML = "";

    courses.forEach(course => {

        container.innerHTML += `
            <div class="course">

                <h3>${course.title}</h3>

                <p>${course.description}</p>

                <p>Price: ₹${course.price}</p>

                <p>Instructor ID: ${course.instructorId}</p>

                <button onclick="purchaseCourse(${course.courseId})">
                    Purchase
                </button>

            </div>
        `;
    });
}


async function purchaseCourse(courseId) {

    if (!userToken) {
        alert("Please login first");
        return;
    }

    const response = await fetch(
        `${API}/user/purchase/${courseId}`,
        {
            method: "POST",

            headers: {
                "Authorization": `Bearer ${userToken}`
            }
        }
    );

    const data = await response.json();

    alert(data.message || "Purchase completed");
}


async function getPurchasedCourses() {

    if (!userToken) {
        alert("Please login first");
        return;
    }

    const response = await fetch(
        `${API}/user/purchased-courses`,
        {
            headers: {
                "Authorization": `Bearer ${userToken}`
            }
        }
    );

    const courses = await response.json();

    const container = document.getElementById("purchased");

    container.innerHTML = "";

    courses.forEach(course => {

        container.innerHTML += `
            <div class="course">

                <h3>${course.title}</h3>

                <p>${course.description}</p>

                <p>${course.content}</p>

            </div>
        `;
    });
}

let adminToken = "";

async function adminLogin() {
    const email = document.getElementById("adminEmail").value;
    const password = document.getElementById("adminPassword").value;

    const response = await fetch(`${API}/admin/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            email: email,
            password: password
        })
    });

    const data = await response.json();

    if (response.ok) {
        adminToken = data.token;
        document.getElementById("message").innerText =
            "Admin login successful!";
    } else {
        document.getElementById("message").innerText =
            data.message || "Admin login failed";
    }
}


async function adminSignup() {
    const name = document.getElementById("adminSignupName").value;
    const email = document.getElementById("adminSignupEmail").value;
    const password = document.getElementById("adminSignupPassword").value;

    const response = await fetch(`${API}/admin/signup`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name: name,
            email: email,
            password: password
        })
    });

    const data = await response.json();

    if (response.ok) {
        document.getElementById("message").innerText =
            "Admin account created successfully!";
    } else {
        document.getElementById("message").innerText =
            data.message || "Admin signup failed";
    }
}






async function createCourse() {
    if (!adminToken) {
        alert("Please login as admin first");
        return;
    }

    const title = document.getElementById("courseTitle").value;
    const description = document.getElementById("courseDescription").value;
    const price = Number(document.getElementById("coursePrice").value);

    const response = await fetch(`${API}/admin/course`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${adminToken}`
        },
        body: JSON.stringify({
            name: title,
            description: description,
            price: price
        })
    });

    const data = await response.json();

    if (response.ok) {
        document.getElementById("message").innerText =
            "Course created successfully!";
    } else {
        document.getElementById("message").innerText =
            data.message || "Course creation failed";
    }
}