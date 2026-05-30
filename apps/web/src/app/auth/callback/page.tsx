"use client"
export default function callbackPage() {
    return <div>
        {document!?.URL ?? ""}
    </div>
}