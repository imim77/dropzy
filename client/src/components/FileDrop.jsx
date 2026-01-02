import React, { useRef, useState } from "react";

import { useNavigate } from 'react-router-dom';

function FileDrop() {
  const uploadRef = useRef(null);
  const [file, setFile] = useState(null);
  const [isDragOver, setIsDragOver] = useState(false);
  const [previewUrl, setPreviewUrl] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleFile = (f) => {
    if (!f) return;
    setFile(f);
    setPreviewUrl(URL.createObjectURL(f));
  };

  const handleChange = (e) => {
    handleFile(e.target.files[0]);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = () => {
    setIsDragOver(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragOver(false);
    handleFile(e.dataTransfer.files[0]);
  };



  

  const handleSubmit = async (e) => {
    e.stopPropagation();

    if (!file) {
      alert("Nema odabrane datoteke");
      return;
    }

    setLoading(true);

    try {
      const formData = new FormData();
      formData.append("file", file);

      const res = await fetch("http://localhost:42069/upload", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        throw new Error("Upload nije uspio");
      }

      const data = await res.json();
      console.log("Upload OK:", data);
      alert("Upload uspješan!");
    } catch (err) {
      console.error(err);
      alert("Greška pri uploadu");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      onClick={() => uploadRef.current.click()}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      style={{
        border: "2px dashed #888",
        padding: "40px",
        textAlign: "center",
        cursor: "pointer",
        background: isDragOver ? "#eef" : "#fafafa",
        position: "relative",
        width: "300px",
        margin: "20px auto",
      }}
    >
      <input
        ref={uploadRef}
        type="file"
        style={{ display: "none" }}
        onChange={handleChange}
      />

      <h2>📂 Upload datoteke</h2>

      {file && <p><strong>{file.name}</strong></p>}

      {previewUrl && (
        <img
          src={previewUrl}
          alt="preview"
          style={{
            width: "120px",
            height: "120px",
            objectFit: "cover",
            marginTop: "10px",
            borderRadius: "6px",
          }}
        />
      )}

      <button
        type="button"
        onClick={handleSubmit}
        disabled={loading}
        style={{
          marginTop: "20px",
          padding: "10px 20px",
          cursor: loading ? "not-allowed" : "pointer",
        }}
      >
        {loading ? "Šaljem..." : "Šalji"}
      </button>


    </div>
  );
}

export default FileDrop;

