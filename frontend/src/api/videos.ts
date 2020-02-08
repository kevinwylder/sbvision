import { session } from './session';

export interface Video {
    id: number
    title: string
    thumbnail: string
    type: number
    format: string
    duration: number
    fps: number
    clips: number
}

interface GetVideosResponse {
    videos: Video[]
    total: number
}

export function getVideos(offset: number, limit: number): Promise<GetVideosResponse> {
    return Promise.resolve({
        videos: [{
            clips: 0,
            duration: 90,
            format: "video/mp4",
            fps: 22,
            id: 2,
            thumbnail: "thumbnail/1U-cgn3cEGA.jpg",
            title: "This is a test video so I can code without the Internet!",
            type: 3,
        }],
        total: 1
    })
    /*return fetch(`/videos?offset=${offset}&count=${limit}`)
    .then(res => {
        if (res.status != 200) {
            return res.text()
            .then(reason => Promise.reject(reason))
        } 
        return res.json()
    }) */
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