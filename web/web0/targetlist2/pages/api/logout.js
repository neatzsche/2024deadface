import { serialize } from 'cookie';

export default function handler(req, res) {
  res.setHeader(
    'Set-Cookie',
    serialize('auth', '', {
      path: '/',
      httpOnly: true,
      expires: new Date(0),
    })
  );
  res.writeHead(302, { Location: '/' });
  res.end();
}