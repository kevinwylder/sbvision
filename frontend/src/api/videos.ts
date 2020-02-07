import { session } from './session';

export interface Video {
    id: number
    title: string
    thumbnail: string
    type: number
    duration: number
    fps: number
    clips: number
}

interface GetVideosResponse {
    videos: Video[]
    total: number
}

export function getVideos(offset: number, limit: number): Promise<GetVideosResponse> {
    return fetch(`/videos?offset=${offset}&count=${limit}`)
    .then(res => {
        if (res.status != 200) {
            return res.text()
            .then(reason => Promise.reject(reason))
        } 
        return res.json()
    }) 
}

export function addVideo(url: string, type: number) : Promise<Video> {
    return fetch(`/videos`, {
        method: "POST",
        body: JSON.stringify({ type, url }),
        headers: {
            Session: session,
        }
    })
    .then(res => {
        if (res.status != 200) {
            return res.text()
            .then(reason => Promise.reject(reason))
        } 
        return res.json()
    }) 
}