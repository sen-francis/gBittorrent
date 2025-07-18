export namespace services {
	
	export class FileUploadResponse {
	    TorrentMetainfo: torrent.TorrentMetainfo;
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new FileUploadResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TorrentMetainfo = this.convertValues(source["TorrentMetainfo"], torrent.TorrentMetainfo);
	        this.Err = source["Err"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TrackerScrapeResponse {
	    Downloaded: number;
	    Seeders: number;
	    Leechers: number;
	    Name: string;
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new TrackerScrapeResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Downloaded = source["Downloaded"];
	        this.Seeders = source["Seeders"];
	        this.Leechers = source["Leechers"];
	        this.Name = source["Name"];
	        this.Err = source["Err"];
	    }
	}

}

export namespace torrent {
	
	export class FileInfo {
	    Length: number;
	    Md5Sum: string;
	    Path: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Length = source["Length"];
	        this.Md5Sum = source["Md5Sum"];
	        this.Path = source["Path"];
	    }
	}
	export class TorrentInfo {
	    PieceLength: number;
	    Pieces: number[];
	    IsPrivate: boolean;
	    DirectoryName: string;
	    FileInfoList: FileInfo[];
	
	    static createFrom(source: any = {}) {
	        return new TorrentInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PieceLength = source["PieceLength"];
	        this.Pieces = source["Pieces"];
	        this.IsPrivate = source["IsPrivate"];
	        this.DirectoryName = source["DirectoryName"];
	        this.FileInfoList = this.convertValues(source["FileInfoList"], FileInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TorrentMetainfo {
	    Info: TorrentInfo;
	    InfoHash: number[];
	    Announce: string;
	    AnnounceList: string[][];
	    CreationDate: number;
	    Comment: string;
	    CreatedBy: string;
	    Encoding: string;
	    Size: number;
	
	    static createFrom(source: any = {}) {
	        return new TorrentMetainfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Info = this.convertValues(source["Info"], TorrentInfo);
	        this.InfoHash = source["InfoHash"];
	        this.Announce = source["Announce"];
	        this.AnnounceList = source["AnnounceList"];
	        this.CreationDate = source["CreationDate"];
	        this.Comment = source["Comment"];
	        this.CreatedBy = source["CreatedBy"];
	        this.Encoding = source["Encoding"];
	        this.Size = source["Size"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

