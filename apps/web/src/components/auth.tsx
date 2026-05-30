"use client";

export function Auth() {
    return <div>
        <button onClick={async () => {
            console.log("Logging in")
        window.location.href = "http://localhost:8000/api/auth/github";
        }}>Github</button>
    </div>
}