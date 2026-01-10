import * as React from "react";
import Table from "@mui/joy/Table";
import Button from "@mui/joy/Button";
import Box from "@mui/joy/Box";
import Modal from "@mui/joy/Modal";
import ModalDialog from "@mui/joy/ModalDialog";
import ModalClose from "@mui/joy/ModalClose";
import Typography from "@mui/joy/Typography";
import Input from "@mui/joy/Input";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

interface Member {
  Name: string;
  PhoneNumber: string;
}

export default function MembersTable({ members = []}: { members?: Member[] }) {
  const [tableData, setTableData] = React.useState<Member[]>(members ?? []);
  const [open, setOpen] = React.useState(false);
  const [editIndex, setEditIndex] = React.useState<number | null>(null);
  const [formData, setFormData] = React.useState<Member>({
    Name: "",
    PhoneNumber: "",
  });
  const [nameDisabled, setNameDisabled] = React.useState(false);

  React.useEffect(() => {
    if (Array.isArray(members)) {
        setTableData(members);
    }
  }, [members]);

  React.useEffect(() => {
    if (tableData.length > 0) {
      saveChanges();
    }
  }, [tableData]);

  const authorize = async (action: Function) => {
    const password = prompt("Bitte Passwort eingeben:");

    if (!password) {
      return;
    }

    const response = await fetch(`${API_BASE}/auth`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ password }),
    });

    if (response.ok) {
      action();
    } else {
      alert("Falsches Passwort.");
    }
  };

  const startEdit = (index: number) => {
    setEditIndex(index);
    setFormData(tableData[index]);
    setOpen(true);
    setNameDisabled(true);
  };

  const saveEdit = () => {
    if (editIndex !== null) {
      const newData = [...tableData];
      newData[editIndex] = formData;
      setTableData(newData);
    } else {
      const newData = [...tableData, formData];
      setTableData(newData);
    }
    setOpen(false);
    setEditIndex(null);
  };

  const deleteMember = (index: number) => {
    setTableData((prevData) => prevData.filter((_, i) => i !== index));
  };

  const startAdd = () => {
    setEditIndex(null);
    setFormData({ Name: "", PhoneNumber: "" });
    setOpen(true);
    setNameDisabled(false);
  };

  const saveChanges = async () => {
    try {
      const response = await fetch(`${API_BASE}/members`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(tableData),
      });

      if (!response.ok) {
        throw new Error("Failed to save changes");
      }
    } catch (error) {
      console.error("Error saving members:", error);
    }
  };

  return (
    <Box sx={{ mb: 2 }}>
      <Table color="primary" variant="soft" aria-label="basic table">
        <thead>
          <tr>
            <th style={{ width: "40%" }}>Name</th>
            <th>Telefonnummer</th>
            <th>Aktionen</th>
          </tr>
        </thead>
        <tbody>
          {tableData.length > 0 && tableData.map((member, index) => (
            <tr key={index}>
              <td>{member.Name}</td>
              <td>{member.PhoneNumber.slice(0, 9) + "..."}</td>
              <td>
                <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
                  <Button
                    color="neutral"
                    onClick={() => authorize(() => startEdit(index))}
                  >
                    Bearbeiten
                  </Button>
                  <Button
                    variant="soft"
                    color="danger"
                    onClick={() => authorize(() => deleteMember(index))}
                  >
                    Löschen
                  </Button>
                </Box>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>

      <Box sx={{ marginTop: 2, display: "flex", gap: 2, flexWrap: "wrap" }}>
        <Button color="primary" onClick={() => authorize(startAdd)}>
          Hinzufügen
        </Button>
      </Box>

      <Modal open={open} onClose={() => setOpen(false)}>
        <ModalDialog>
          <ModalClose />
          <Typography level="h4" sx={{ mb: 2 }}>
            Mitglied
          </Typography>
          <Box sx={{ display: "flex", flexDirection: "column", gap: 2 }}>
            <Input
              disabled={nameDisabled}
              placeholder="Name"
              value={formData.Name}
              onChange={(e) =>
                setFormData({ ...formData, Name: e.target.value })
              }
            />
            <Input
              placeholder="Telefonnummer"
              value={formData.PhoneNumber}
              onChange={(e) =>
                setFormData({ ...formData, PhoneNumber: e.target.value })
              }
            />
            <Button color="success" onClick={saveEdit}>
              Speichern
            </Button>
          </Box>
        </ModalDialog>
      </Modal>
    </Box>
  );
}
