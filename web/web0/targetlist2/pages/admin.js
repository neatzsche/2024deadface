// pages/admin.js

import { redirect } from "next/dist/server/api-utils";

export default function Admin({ username }) {
    return (
      <div>
        <h1>Admin Page</h1>
        <p>Welcome, {username}!</p>
        <a href="/api/logout">Logout</a>
      </div>
    );
  }
  
  export async function getServerSideProps({ req }) {
    const { auth } = req.cookies;
  
    if (!auth) {
      return {
        redirect: {
          destination: '/login',
          permanent: false,
        },
      };
    }
    try{
    const username = auth.split(":")[0]
    const key = auth.split(":")[1]

    console.log(key)

    if (key != 'secretkey729374jhjh98279'){
        return {
            props: { username: 'bad login' },
        }
    }

    if (username == 'admin'){
        return{
        props: { username: 'flag{sh@r3d-s3cr3t-in-D@-c00k13}' },
    };

    }
  
    return {
      props: { username: username },
    };
} catch {
    return {
        props: { username: 'failed' },
      };
}
  }
  