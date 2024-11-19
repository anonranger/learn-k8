<script>
  import { onMount } from 'svelte';

  let users = [];
  let name = "";
  let email = "";
  let age = null;
  let mobile_number = "";

  const backendUrl = import.meta.env.VITE_BACKEND_URL;
  console.log("🍷 ~ backendUrl:", backendUrl);

  onMount(() => {
    fetchUsers();
  });

  async function addUser() {
    if (name || email || age || mobile_number) {
      try {
        const response = await fetch(`${backendUrl}/api/record`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ name, email, age, mobile_number }),
        });

        if (!response.ok) {
          alert("Failed to add user");
          throw new Error("Failed to add user");
        }
        const data = await response.json();
        users = [...users, data.data];
      } catch (error) {
        console.error("Error:", error);
      } finally {
        name = "";
        email = "";
        age = null;
        mobile_number = "";
      }
    }
  }

  async function fetchUsers() {
    try {
      const response = await fetch(`${backendUrl}/api/record`);
      if (!response.ok) {
        throw new Error("Failed to fetch users");
      }

      users = await response.json();
      users = users || [];
    } catch (error) {
      alert("Failed to fetch users");
      console.error("Error:", error);
    }
  }
</script>

<div class="container">
  <div class="form-container">
    <h2>User Details Form</h2>
    <form on:submit|preventDefault={addUser}>
      <div>
        <label for="name">Name:</label>
        <input
          type="text"
          id="name"
          bind:value={name}
        />
      </div>
      <div>
        <label for="email">Email:</label>
        <input
          type="email"
          id="email"
          bind:value={email}
        />
      </div>
      <div>
        <label for="age">Age:</label>
        <input
          type="number"
          id="age"
          bind:value={age}
        />
      </div>
      <div>
        <label for="mobile_number">Mobile Number:</label>
        <input
          type="text"
          id="mobile_number"
          bind:value={mobile_number}
        />
      </div>
      <button type="submit">Add User</button>
    </form>
  </div>

  <div class="table-container">
    <h2>Inserted Users</h2>
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Email</th>
          <th>Age</th>
          <th>Mobile Number</th>
        </tr>
      </thead>
      <tbody>
        {#if !users || users.length === 0}
          <tr>
            <td colspan="4">No users available</td>
          </tr>
        {:else}
          {#each users as user}
          <tr>
            <td>{user.name}</td>
            <td>{user.email}</td>
            <td>{user.age}</td>
            <td>{user.mobile_number}</td>
          </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
