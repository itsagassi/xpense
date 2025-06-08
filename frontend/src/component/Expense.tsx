import React, { useState } from "react";

interface Expense {
  id: number;
  title: string;
  description: string;
  category: string;
  amount: number;
  date: Date;
}

interface Props {
  expenses: Expense[];
  fetchExpenses: () => void;
  categories: string[];
}

export const ExpensesTable = ({ expenses, fetchExpenses, categories }: Props) => {
  const [filter, setFilter] = useState("All");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editForm, setEditForm] = useState<Omit<Expense, "id">>({
    title: "",
    description: "",
    category: "",
    amount: 0,
    date: new Date(),
  });

  const filteredExpenses =
    filter === "All"
      ? expenses
      : expenses.filter((exp) => exp.category === filter);

  const startEdit = (expense: Expense) => {
    setEditingId(expense.id);
    setEditForm({
      title: expense.title,
      description: expense.description,
      category: expense.category,
      amount: expense.amount,
      date: new Date(expense.date),
    });
  };

  const saveEdit = async () => {
    if (editingId === null) return;

    try {
      const res = await fetch(`http://localhost:5000/api/v1/expenses/${editingId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: localStorage.getItem("token") || "",
        },
        body: JSON.stringify(editForm),
      });

      if (!res.ok) throw new Error("Failed to update expense");

      setEditingId(null);
      fetchExpenses(); // refetch after edit
    } catch (error) {
      console.error("Error updating expense:", error);
    }
  };

  const cancelEdit = () => {
    setEditingId(null);
  };

  const deleteExpense = async (id: number) => {
    if (!window.confirm("You sure you wanna delete this expense?")) return;

    try {
      const res = await fetch(`http://localhost:5000/api/v1/expenses/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: localStorage.getItem("token") || "",
        },
      });

      if (!res.ok) throw new Error("Failed to delete expense");

      fetchExpenses(); // refetch after delete
    } catch (error) {
      console.error("Error deleting expense:", error);
    }
  };

  const formatDateForInput = (date: Date) => date.toISOString().slice(0, 10);

  return (
    <div style={styles.container}>
      <h2>Expenses</h2>

      <div style={styles.filterContainer}>
        <label>
          Filter by Category:{" "}
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            style={styles.select}
          >
            {categories.map((cat) => (
              <option key={cat} value={cat}>
                {cat}
              </option>
            ))}
          </select>
        </label>
      </div>

      <table style={styles.table}>
        <thead>
          <tr>
            <th>Title</th>
            <th>Description</th>
            <th>Category</th>
            <th>Amount (Rp)</th>
            <th>Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {filteredExpenses.map((exp) =>
            editingId === exp.id ? (
              <tr key={exp.id} style={styles.editRow}>
                <td style={styles.tableValue}>
                  <input
                    type="text"
                    value={editForm.title}
                    onChange={(e) =>
                      setEditForm((f) => ({ ...f, title: e.target.value }))
                    }
                    style={styles.input}
                  />
                </td>
                <td style={styles.tableValue}>
                  <input
                    type="text"
                    value={editForm.description}
                    onChange={(e) =>
                      setEditForm((f) => ({ ...f, description: e.target.value }))
                    }
                    style={styles.input}
                  />
                </td>
                <td style={styles.tableValue}>
                  <select
                    value={editForm.category}
                    onChange={(e) =>
                      setEditForm((f) => ({ ...f, category: e.target.value }))
                    }
                    style={styles.select}
                  >
                    {categories
                      .filter((c) => c !== "All")
                      .map((cat) => (
                        <option key={cat} value={cat}>
                          {cat}
                        </option>
                      ))}
                  </select>
                </td>
                <td style={styles.tableValue}>
                  <input
                    type="number"
                    value={editForm.amount}
                    onChange={(e) =>
                      setEditForm((f) => ({ ...f, amount: +e.target.value }))
                    }
                    style={styles.input}
                    min={0}
                  />
                </td>
                <td style={styles.tableValue}>
                  <input
                    type="date"
                    value={formatDateForInput(editForm.date)}
                    onChange={(e) =>
                      setEditForm((f) => ({ ...f, date: new Date(e.target.value) }))
                    }
                    style={styles.input}
                  />
                </td>
                <td style={styles.tableValue}>
                  <button onClick={saveEdit} style={styles.saveBtn}>
                    Save
                  </button>
                  <button onClick={cancelEdit} style={styles.cancelBtn}>
                    Cancel
                  </button>
                </td>
              </tr>
            ) : (
              <tr key={exp.id}>
                <td style={styles.tableValue}>{exp.title}</td>
                <td style={styles.tableValue}>{exp.description}</td>
                <td style={styles.tableValue}>{exp.category}</td>
                <td style={styles.tableValue}>{exp.amount}</td>
                <td style={styles.tableValue}>
                  {new Date(exp.date).toLocaleDateString()}
                </td>
                <td style={styles.tableValue}>
                  <button onClick={() => startEdit(exp)} style={styles.actionBtn}>
                    Edit
                  </button>
                  <button
                    onClick={() => deleteExpense(exp.id)}
                    style={{ ...styles.actionBtn, ...styles.deleteBtn }}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            )
          )}
        </tbody>
      </table>
    </div>
  );
};

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    display: "flex",
    flexDirection: "column",
    backgroundColor: "white",
    boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
    borderRadius: "0.5rem",
    padding: "1rem",
    marginBottom: "1rem",
    border: "1px solid #e5e7eb",
  },
  filterContainer: {
    marginBottom: "1rem",
  },
  select: {
    padding: "0.3rem 0.6rem",
    borderRadius: "6px",
    border: "1px solid #ccc",
    fontSize: "1rem",
  },
  table: {
    width: "100%",
    borderCollapse: "collapse",
  },
  input: {
    width: "100%",
    padding: "0.3rem 0.5rem",
    fontSize: "1rem",
    borderRadius: "6px",
    border: "1px solid #ccc",
  },
  editRow: {
    backgroundColor: "#e8f0fe",
  },
  actionBtn: {
    marginRight: "0.5rem",
    padding: "0.3rem 0.8rem",
    fontSize: "0.9rem",
    borderRadius: "6px",
    border: "none",
    cursor: "pointer",
    backgroundColor: "#2196f3",
    color: "white",
    transition: "background-color 0.3s ease",
  },
  deleteBtn: {
    backgroundColor: "#f44336",
  },
  saveBtn: {
    padding: "0.3rem 0.8rem",
    marginRight: "0.5rem",
    backgroundColor: "#4caf50",
    border: "none",
    borderRadius: "6px",
    color: "white",
    cursor: "pointer",
  },
  cancelBtn: {
    padding: "0.3rem 0.8rem",
    backgroundColor: "#9e9e9e",
    border: "none",
    borderRadius: "6px",
    color: "white",
    cursor: "pointer",
  },
  tableValue: {
    padding: "0.5rem",
    textAlign: "center",
    borderBottom: "1px solid #ddd",
  },
};
