import React, { useState } from "react";
import { createClient } from "@supabase/supabase-js";
import { useNavigate } from "react-router-dom";

// Initialize Supabase
const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL!,
  import.meta.env.VITE_SUPABASE_ANON_KEY!
);

export const AuthPage = () => {
  const [isRegister, setIsRegister] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [showResend, setShowResend] = useState(false);
  const navigate = useNavigate(); // 👈 use this to redirect

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    setShowResend(false);

    let supabaseError = null;
    let session = null;

    if (isRegister) {
      const { error } = await supabase.auth.signUp({ email, password });
      supabaseError = error;
    } else {
      const { data, error } = await supabase.auth.signInWithPassword({ email, password });
      supabaseError = error;
      session = data?.session;
    }

    setLoading(false);

    if (supabaseError) {
      if (supabaseError.message === "Email not confirmed") {
        setError("Please confirm your email first. Didn't get it? Resend below.");
        setShowResend(true);
      } else {
        setError(supabaseError.message);
        setShowResend(false);
      }
    } else {
      setError(null);
      if (isRegister) {
        alert("Registration successful! Please check your email to confirm.");
      } else {
        localStorage.setItem("token", "Bearer " + (session?.access_token ?? ""));
        navigate("/dashboard");
      }
    }
  };

  const handleResendConfirmation = async () => {
    setLoading(true);
    setError(null);
    const { error } = await supabase.auth.signInWithOtp({ email });
    setLoading(false);

    if (error) {
      setError(error.message);
    } else {
      alert("New confirmation email sent! Check your inbox.");
      setShowResend(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-box">
        <h1>{isRegister ? "Register for Xpense" : "Login to Xpense"}</h1>

        <form onSubmit={handleSubmit} className="auth-form">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="auth-input"
          />

          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="auth-input"
          />

          {error && <p className="auth-error">{error}</p>}

          <button type="submit" disabled={loading} className="auth-button">
            {loading ? "Please wait..." : isRegister ? "Register" : "Login"}
          </button>
        </form>

        {showResend && (
          <button
            onClick={handleResendConfirmation}
            disabled={loading}
            className="resend-button"
          >
            Resend Confirmation Email
          </button>
        )}

        <p className="auth-switch">
          {isRegister ? "Already have an account?" : "Don't have an account?"}{" "}
          <button
            onClick={() => {
              setIsRegister(!isRegister);
              setError(null);
              setShowResend(false);
            }}
            className="switch-button"
          >
            {isRegister ? "Login" : "Register"}
          </button>
        </p>
      </div>
    </div>
  );
};
