import React from "react";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
} from "recharts";

interface ExpenseChartsProps {
  pieData: { name: string; value: number }[];
  barDataMonth: { name: string; total: number }[];
  barDataWeek: { name: string; total: number }[];
}

const COLORS = [
  "#1f78b4", // deep blue
  "#33a02c", // forest green
  "#e31a1c", // crimson red
  "#ff7f00", // burnt orange
  "#6a3d9a", // dark purple
  "#b15928", // dark gold
  "#a6cee3", // steel blue
  "#b2df8a", // olive green
  "#fb9a99", // dusty pink
  "#fdbf6f", // warm amber
  "#cab2d6", // muted lilac
  "#ffff99", // pale yellow
];



export const Summary: React.FC<ExpenseChartsProps> = ({
  pieData,
  barDataMonth,
  barDataWeek,
}) => {
  return (
    <>
      <div style={styles.container}>
        <div style={styles.chartBox}>
          <h3 style={styles.chartTitle}>Expenses by Category</h3>
          <PieChart width={300} height={300}>
            <Pie
              data={pieData}
              dataKey="value"
              nameKey="name"
              cx="50%"
              cy="50%"
              outerRadius={100}
              label
            >
              {pieData.map((_entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={COLORS[index % COLORS.length]}
                />
              ))}
            </Pie>
            <Tooltip />
            <Legend verticalAlign="bottom" height={36} />
          </PieChart>
        </div>

        <div style={styles.chartBox}>
          <h3 style={styles.chartTitle}>Expense By Caregory</h3>
          <BarChart width={500} height={300} data={pieData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" tick={{ fill: "#666" }} />
            <YAxis tick={{ fill: "#666" }} />
            <Tooltip />
            <Legend verticalAlign="bottom" height={36} />
            <Bar dataKey="value" fill="#1f78b4" radius={[4, 4, 0, 0]} />
          </BarChart>
        </div>
      </div>
      <div style={styles.container}>
        <div style={styles.chartBox}>
          <h3 style={styles.chartTitle}>Expenses Over Time (Month)</h3>
          <BarChart width={500} height={300} data={barDataMonth}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" tick={{ fill: "#666" }} />
            <YAxis tick={{ fill: "#666" }} />
            <Tooltip />
            <Legend verticalAlign="bottom" height={36} />
            <Bar dataKey="total" fill="#33a02c" radius={[4, 4, 0, 0]} />
          </BarChart>
        </div>

        <div style={styles.chartBox}>
          <h3 style={styles.chartTitle}>Expenses Over Time (Week)</h3>
          <BarChart width={500} height={300} data={barDataWeek}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" tick={{ fill: "#666" }} />
            <YAxis tick={{ fill: "#666" }} />
            <Tooltip />
            <Legend verticalAlign="bottom" height={36} />
            <Bar dataKey="total" fill="#e31a1c" radius={[4, 4, 0, 0]} />
          </BarChart>
        </div>
      </div>
    </>
  );
};

const styles: { [key: string]: React.CSSProperties } = {
  container: {
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    backgroundColor: "white",
    boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
    borderRadius: "0.5rem",
    padding: "1rem",
    marginBottom: "1rem",
    border: "1px solid #e5e7eb",
  },
  chartBox: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    background: "#f0f4f8",
    padding: "1.5rem",
    borderRadius: "12px",
    boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
    minWidth: "450px",
  },
  chartTitle: {
    marginBottom: "1rem",
    fontFamily: "'Segoe UI', Tahoma, Geneva, Verdana, sans-serif",
    color: "#333",
  },
};
