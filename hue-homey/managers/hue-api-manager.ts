import fetch from 'node-fetch';

const api_url = 'https://10.0.6.106/clip/v2';
const username = '4d1CmbqcaWReBdlaNmjk5cstX68xT-hDG5MxjdCl';

process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = "0";

export async function activateScene(id: string): Promise<boolean> {
  const response = await fetch(`${api_url}/resource/scene/${id}`, {
    method: 'PUT',
    headers: {
      'Accept': 'application/json',
      'hue-application-key': username
    },
    body: JSON.stringify({
      recall: {
        action: 'active'
      }
    })
  });

  if (!response.ok) {
    throw new Error(response.statusText);
  }

  const json = await response.json();

  return json.error && json.error.length === 0;
}

export async function activateSmartScene(id: string): Promise<boolean> {
  const response = await fetch(`${api_url}/resource/smart_scene/${id}`, {
    method: 'PUT',
    headers: {
      'Accept': 'application/json',
      'hue-application-key': username
    },
    body: JSON.stringify({
      recall: {
        action: 'activate'
      }
    })
  });

  if (!response.ok) {
    throw new Error(response.statusText);
  }

  const json = await response.json();

  return json.error && json.error.length === 0;
}

export async function getSceneIdByName(name: string): Promise<string> {
  const response = await fetch(`${api_url}/resource/scene`, {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
      'hue-application-key': username
    }
  });

  if (!response.ok) {
    throw new Error(response.statusText);
  }

  const json = await response.json();
  const scene = json.data.find((element: any) => element.metadata.name === name);

  if (scene) {
    return scene.id;
  } else {
    throw new Error('Smart Scene not found with name: ' + name);
  }
}

export async function getSmartSceneIdByName(name: string): Promise<string> {
  const response = await fetch(`${api_url}/resource/smart_scene`, {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
      'hue-application-key': username
    }
  });

  if (!response.ok) {
    throw new Error(response.statusText);
  }

  const json = await response.json();
  const scene = json.data.find((element: any) => element.metadata.name === name);

  if (scene) {
    return scene.id;
  } else {
    throw new Error('Smart Scene not found with name: ' + name);
  }
}

export default activateScene
