import { useState } from "react";

interface Expense {
  title: string;
  description: string;
  category: string;
  amount: number;
  date: Date;
}

interface CreateExpenseProps {
  onExpenseCreated: () => void;
}

export const CreateExpense = ({ onExpenseCreated }: CreateExpenseProps) => {
  const [form, setForm] = useState<Expense>({
    title: "",
    description: "",
    amount: 0,
    category: "",
    date: new Date(),
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    if (name === "amount") {
      setForm({ ...form, [name]: +value });
    } else if (name === "date") {
      setForm({ ...form, date: new Date(value) });
    } else {
      setForm({ ...form, [name]: value });
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:5000/api/v1/expenses", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: localStorage.getItem("token") || "",
        },
        body: JSON.stringify(form),
      });

      if (!res.ok) throw new Error("Failed to add expense");

      setForm({
        title: "",
        description: "",
        amount: 0,
        category: "",
        date: new Date(),
      });

      onExpenseCreated();
    } catch (error) {
      console.error("Error adding expense:", error);
    }
  };

  const formatDateForInput = (date: Date) => date.toISOString().slice(0, 10);

  return (
    <form onSubmit={handleSubmit} style={styles.card}>
      <h2 style={styles.heading}>Add New Expense</h2>
      <input
        name="title"
        type="text"
        placeholder="Title"
        value={form.title}
        onChange={handleChange}
        style={styles.input}
        required
      />
      <input
        name="description"
        type="text"
        placeholder="Description"
        value={form.description}
        onChange={handleChange}
        style={styles.input}
        required
      />
      <input
        name="amount"
        type="number"
        placeholder="Amount"
        value={form.amount}
        onChange={handleChange}
        style={styles.input}
        required
      />
      <input
        name="category"
        type="text"
        placeholder="Category"
        value={form.category}
        onChange={handleChange}
        style={styles.input}
        required
      />
      <input
        name="date"
        type="date"
        value={formatDateForInput(form.date)}
        onChange={handleChange}
        style={styles.input}
        required
      />
      <button type="submit" style={styles.button}>
        Add Expense
      </button>
    </form>
  );
};

const styles: { [key: string]: React.CSSProperties } = {
  card: {
    display: "flex",
    flexDirection: "column",
    backgroundColor: "white",
    boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
    borderRadius: "0.5rem",
    padding: "1rem",
    marginBottom: "1rem",
    border: "1px solid #e5e7eb",
  },
  heading: {
    fontSize: "1.125rem",
    fontWeight: "600",
    marginBottom: "0.5rem",
  },
  input: {
    width: "100%",
    marginBottom: "0.5rem",
    padding: "0.5rem 1rem",
    border: "1px solid #d1d5db",
    borderRadius: "0.375rem",
    boxSizing: "border-box",
  },
  button: {
    width: "100%",
    backgroundColor: "#3b82f6",
    color: "white",
    padding: "0.5rem 1rem",
    borderRadius: "0.375rem",
    border: "none",
    cursor: "pointer",
  },
};
